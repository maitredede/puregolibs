package evdi

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/strings"
	"golang.org/x/sys/unix"
)

type evdiHandle unsafe.Pointer

func AddDevice() int {
	initLib()
	return int(libEvdiAddDevice())
}

type Device struct {
	h evdiHandle

	// lck     sync.Mutex
	// evtChan chan Events
	nextBufferID atomic.Int32
	buffersMap   map[int32]bufferData
	buffersArr   []bufferData

	bufferToUpdate int
}

type bufferData struct {
	buffer *evdiBuffer
	data   []byte
}

func OpenDevice(device int) (*Device, error) {
	initLib()

	h := libEvdiOpen(int32(device))
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		h:          h,
		buffersMap: make(map[int32]bufferData),
		buffersArr: make([]bufferData, 0, 2),
	}
	return d, nil
}

func OpenAttachedToNone() (*Device, error) {
	initLib()

	h := libEvdiOpenAttachedToFixed(nil, 0)
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		h:          h,
		buffersMap: make(map[int32]bufferData),
		buffersArr: make([]bufferData, 0, 2),
	}
	return d, nil
}

func OpenAttachedTo(sysfsParent string) (*Device, error) {
	initLib()

	cParent, l := strings.CStringL(sysfsParent)
	h := libEvdiOpenAttachedToFixed(unsafe.Pointer(cParent), uint32(l))
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		h:          h,
		buffersMap: make(map[int32]bufferData),
		buffersArr: make([]bufferData, 0, 2),
	}
	return d, nil
}

func (d *Device) Close() error {
	initLib()
	if d.h == nil {
		return ErrDeviceIsClosed
	}

	libEvdiClose(d.h)
	d.h = nil

	return nil
}

func (d *Device) EnableCursorEvents(enabled bool) error {
	initLib()
	if d.h == nil {
		return ErrDeviceIsClosed
	}
	libEvdiEnableCursorEvents(d.h, enabled)
	return nil
}

func (d *Device) DummyEDID() []byte {
	return EDIDv1_1280x800[:]
}

func (d *Device) RunDummy(close context.Context) error {
	initLib()
	if d.h == nil {
		return ErrDeviceIsClosed
	}

	edid := d.DummyEDID()
	eventsHandler := Events{
		ModeChanged: func(mode EvdiMode, userdata unsafe.Pointer) {
			slog.Info(fmt.Sprintf("mode: %+v", mode))
			for bid, _ /* buff */ := range d.buffersMap {
				libevdiUnregisterBuffer(d.h, int32(bid))
				delete(d.buffersMap, bid)
			}
			d.buffersArr = make([]bufferData, 0, 2)

			const bufferCount = 2
			for i := 0; i < bufferCount; i++ {
				id := d.nextBufferID.Add(1)
				buffer := &evdiBuffer{
					id:     id,
					width:  mode.Width,
					height: mode.Height,
					stride: mode.BitsPerPixel / 8 * mode.Height,
				}
				data := make([]byte, buffer.height*buffer.stride)
				buffer.buffer = &data[0]

				libevdiRegisterBuffer(d.h, buffer)
				b := bufferData{
					data:   data,
					buffer: buffer,
				}
				d.buffersMap[buffer.id] = b
				d.buffersArr = append(d.buffersArr, b)
			}
		},

		Dpms: func(dpmsMode int32, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("dpms: 0x%x", dpmsMode))
		},
		UpdateReady: func(bufferToUpdate int32, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("updateReady: %d", bufferToUpdate))
		},
		CrtcState: func(state int32, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("crtc: 0x%x", state))
		},
		CursorSet: func(cursorSet CursorSet, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("cursorSet: %+v", cursorSet))
		},
		CursorMove: func(cursorMove CursorMove, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("cursorMove: %+v", cursorMove))
		},
		DdcCiData: func(data DdcciData, userData unsafe.Pointer) {
			slog.Info(fmt.Sprintf("ddciData: %+v", data))
		},
	}

	nativeEvts, dispose := buildNativeEvents(eventsHandler)
	defer dispose()
	slog.Debug("evdi: connecting")
	maxUint32 := uint32(0xFFFFFFFF)
	libEvdiConnect(d.h, &edid[0], uint32(len(edid)), maxUint32)
	slog.Debug("evdi: connected")
	defer libEvdiDisconnect(d.h)

	newSelectable := libEvdiGetEventReady(d.h)
	slog.Debug(fmt.Sprintf("evdi: eventGetReady: %v", newSelectable))

	epollFd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(fmt.Errorf("failed to create epoll: %w", err))
	}
	defer unix.Close(epollFd)

	// Configurer epoll pour surveiller les événements READ sur ce fd
	event := unix.EpollEvent{
		Events: unix.EPOLLIN, // Équivalent de ev::READ
		Fd:     int32(newSelectable),
	}
	if err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, int(newSelectable), &event); err != nil {
		panic(fmt.Errorf("failed to add fd to epoll: %w", err))
	}
	defer func() {
		// Retirer le fd d'epoll
		if err := unix.EpollCtl(epollFd, unix.EPOLL_CTL_DEL, int(newSelectable), nil); err != nil {
			panic(fmt.Errorf("failed to remove fd from epoll: %w", err))
		}
	}()

	//wait for read
	events := make([]unix.EpollEvent, 10)
	for {
		select {
		case <-close.Done():
			return close.Err()
		default:
		}

		// Attendre les événements (timeout de 100ms pour vérifier done)
		n, err := unix.EpollWait(epollFd, events, 100)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			fmt.Printf("epoll wait error: %v\n", err)
			continue
		}

		// Traiter chaque événement
		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			slog.Debug(fmt.Sprintf("using: event i=%v fd=%v", i, fd))
		}

		// }

		// for {
		// 	select {
		// 	case <-close.Done():
		// 		return close.Err()
		// 	default:
		// 	}

		libEvdiHandleEvents(d.h, nativeEvts)

		// slog.Debug(".")
		if len(d.buffersArr) > 0 {
			bd := d.buffersArr[d.bufferToUpdate]
			d.bufferToUpdate = (d.bufferToUpdate + 1) % len(d.buffersArr)
			update := libevdiRequestUpdate(d.h, bd.buffer.id)
			if update {
				var numRects int32
				var rects *evdiRect
				libevdiGrabPixels(d.h, &rects, &numRects)
			}
		}
	}
}

// func (d *Device) runConnect(close context.Context, edid []byte, skuAreaLimit uint32) (<-chan Events, error) {
// 	initLib()
// 	if d.h == nil {
// 		return nil, ErrDeviceIsClosed
// 	}

// 	d.lck.Lock()
// 	defer d.lck.Unlock()

// 	if d.evtChan != nil {
// 		return nil, fmt.Errorf("already handling events")
// 	}

// 	d.evtChan = make(chan Events, 1)

// 	edidPtr := (&edid[0])
// 	edidLen := uint32(len(edid))
// 	libEvdiConnect(d.h, edidPtr, edidLen, skuAreaLimit)
// 	go d.runLoop(close, d.evtChan)
// 	return d.evtChan, nil
// }

// func (d *Device) runLoop(close context.Context, evtChan chan<- Events) {
// 	defer libEvdiDisconnect(d.h)
// 	defer func() {
// 		d.lck.Lock()
// 		defer d.lck.Unlock()
// 		d.evtChan = nil
// 	}()
// 	for {
// 		select {
// 		case <-close.Done():
// 			return
// 		default:
// 		}

// 	}
// }

// func (d *Device) Connect(edid []byte, skuAreaLimit uint32) error {
// 	initLib()
// 	if d.h == nil {
// 		return ErrDeviceIsClosed
// 	}
// }

// func (d *Device) Connect2(edid []byte, pixelAreaLimit uint32, pixelPerSecondLimit uint32) error {
// 	initLib()
// 	if d.h == nil {
// 		return ErrDeviceIsClosed
// 	}
// }

// func (d *Device) Disconnect() error {
// 	initLib()
// 	if d.h == nil {
// 		return ErrDeviceIsClosed
// 	}
// }

// func (d *Device) RegisterBuffer(width, height, stride, rectCount int) (*Buffer, error) {
// }

func buildNativeEvents(e Events) (*evdiEventContext, func()) {
	disposables := make([]func(), 0, 10)
	dispose := func() {
		for _, d := range disposables {
			d()
		}
	}

	ne := &evdiEventContext{}
	//DPMS
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})
		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			modePtr := argArr[0]
			userDataPtr := argArr[1]

			mode := *(*int32)(modePtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.Dpms(mode, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.dpmsHandler = callback
	}
	// ModeChanged
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})

		cifEvdiMode := ffi.NewType(
			&ffi.TypeSint32, // width
			&ffi.TypeSint32, // height
			&ffi.TypeSint32, // refresh_rate
			&ffi.TypeSint32, // bits_per_pixel
			&ffi.TypeUint32, // pixel_format
		)

		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &cifEvdiMode, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			modePtr := argArr[0]
			userDataPtr := argArr[1]

			mode := *(*EvdiMode)(modePtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.ModeChanged(mode, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.modeChangedHandler = callback
	}
	//updateReady
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})
		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			bufferPtr := argArr[0]
			userDataPtr := argArr[1]

			buffer := *(*int32)(bufferPtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.UpdateReady(buffer, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.updateReadyHandler = callback
	}
	//crtcState
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})
		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			statePtr := argArr[0]
			userDataPtr := argArr[1]

			state := *(*int32)(statePtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.CrtcState(state, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.crtcStateHandler = callback
	}
	//cursorSet
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})

		cifEvdiCursorSet := ffi.NewType(
			&ffi.TypeSint32,  // hotX
			&ffi.TypeSint32,  // hotY
			&ffi.TypeUint32,  // width
			&ffi.TypeUint32,  // height
			&ffi.TypeUint8,   // enabled
			&ffi.TypeUint32,  // bufferLength
			&ffi.TypePointer, // buffer
			&ffi.TypeUint32,  // pixelFormat
			&ffi.TypeUint32,  // stride
		)

		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &cifEvdiCursorSet, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			cursorSetPtr := argArr[0]
			userDataPtr := argArr[1]

			cursorSet := *(*CursorSet)(cursorSetPtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.CursorSet(cursorSet, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.cursorSetHandler = callback
	}
	//cursor move
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})

		cifEvdiCursorMove := ffi.NewType(
			&ffi.TypeSint32, // x
			&ffi.TypeSint32, // y
		)

		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &cifEvdiCursorMove, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			cursorMovePtr := argArr[0]
			userDataPtr := argArr[1]

			cursorMove := *(*CursorMove)(cursorMovePtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			e.CursorMove(cursorMove, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.dpmsHandler = callback
	}
	//ddcci
	{
		var callback unsafe.Pointer
		closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
		if closure == nil {
			panic("closureAlloc")
		}
		disposables = append(disposables, func() {
			ffi.ClosureFree(closure)
		})

		cifDcciData := ffi.NewType(
			&ffi.TypeUint16,  // address
			&ffi.TypeUint16,  // flag
			&ffi.TypeUint32,  // buffer_length
			&ffi.TypePointer, // buffer
		)

		var cifCallback ffi.Cif
		if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 2, &ffi.TypeVoid, &cifDcciData, &ffi.TypePointer); status != ffi.OK {
			panic("cif")
		}
		fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, ffiUserData unsafe.Pointer) uintptr {
			argArr := unsafe.Slice(args, cif.NArgs)
			dataPtr := argArr[0]
			userDataPtr := argArr[1]

			nativeData := *(*evdiDdcciData)(dataPtr)
			userData := *(*unsafe.Pointer)(userDataPtr)

			data := DdcciData{
				Address: nativeData.address,
				Flags:   nativeData.flags,
				Buffer:  unsafe.Slice(nativeData.buffer, nativeData.bufferLength),
			}

			e.DdcCiData(data, userData)
			return 0
		})
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
		ne.ddcciDataHandler = callback
	}

	return ne, dispose
}

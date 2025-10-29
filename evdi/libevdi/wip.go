package libevdi

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/maitredede/puregolibs/evdi"
)

var EDIDv1_1280x800 = evdi.EDIDv1_1280x800

func OpenAttachedToNone() (*Handle, error) {
	return OpenAttachedTo("")
}

func (h *Handle) DummyEDID() []byte {
	return EDIDv1_1280x800[:]
}

func (h *Handle) RunDummy(ctx context.Context) error {
	// eventsHandler := Events{
	// 	ModeChanged: func(mode EvdiMode, userdata unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("mode: %+v", mode))
	// 		for bid, _ /* buff */ := range d.buffersMap {
	// 			libevdiUnregisterBuffer(d.h, int32(bid))
	// 			delete(d.buffersMap, bid)
	// 		}
	// 		d.buffersArr = make([]bufferData, 0, 2)

	// 		const bufferCount = 2
	// 		for i := 0; i < bufferCount; i++ {
	// 			id := d.nextBufferID.Add(1)
	// 			buffer := &evdiBuffer{
	// 				id:     id,
	// 				width:  mode.Width,
	// 				height: mode.Height,
	// 				stride: mode.BitsPerPixel / 8 * mode.Height,
	// 			}
	// 			data := make([]byte, buffer.height*buffer.stride)
	// 			buffer.buffer = &data[0]

	// 			libevdiRegisterBuffer(d.h, buffer)
	// 			b := bufferData{
	// 				data:   data,
	// 				buffer: buffer,
	// 			}
	// 			d.buffersMap[buffer.id] = b
	// 			d.buffersArr = append(d.buffersArr, b)
	// 		}
	// 	},

	// 	Dpms: func(dpmsMode int32, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("dpms: 0x%x", dpmsMode))
	// 	},
	// 	UpdateReady: func(bufferToUpdate int32, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("updateReady: %d", bufferToUpdate))
	// 	},
	// 	CrtcState: func(state int32, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("crtc: 0x%x", state))
	// 	},
	// 	CursorSet: func(cursorSet CursorSet, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("cursorSet: %+v", cursorSet))
	// 	},
	// 	CursorMove: func(cursorMove CursorMove, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("cursorMove: %+v", cursorMove))
	// 	},
	// 	DdcCiData: func(data DdcciData, userData unsafe.Pointer) {
	// 		slog.Info(fmt.Sprintf("ddciData: %+v", data))
	// 	},
	// }

	// nativeEvts, dispose := buildNativeEvents(eventsHandler, nil)
	// defer dispose()
	slog.Debug("dummy: connecting")
	// maxUint32 := uint32(0xFFFFFFFF)

	skuLimit := uint32(1280 * 800)
	h.Connect(EDIDv1_1280x800, skuLimit)
	slog.Debug("dummy: connected")
	defer h.Disconnect()

	newSelectable := h.GetEventReady()
	slog.Debug(fmt.Sprintf("dummy: eventGetReady: %v", newSelectable))

	// // Configurer epoll pour surveiller les événements READ sur ce fd
	// event := unix.PollFd{
	// 	Events: unix.EPOLLIN, // Équivalent de ev::READ
	// 	Fd:     int32(newSelectable),
	// }
	// events := []unix.PollFd{
	// 	event,
	// }

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	default:
	// 	}

	// 	// Attendre les événements (timeout de 100ms pour vérifier done)
	// 	n, err := unix.Poll(events, 1000)
	// 	if err != nil {
	// 		if err == unix.EINTR {
	// 			continue
	// 		}
	// 		fmt.Printf("epoll wait error: %v\n", err)
	// 		continue
	// 	}

	// 	// Traiter chaque événement
	// 	for i := 0; i < n; i++ {
	// 		fd := int(events[i].Fd)
	// 		slog.Debug(fmt.Sprintf("using: event i=%v fd=%v", i, fd))
	// 	}

	// 	// libEvdiHandleEvents(d.h, nativeEvts)

	// 	// // slog.Debug(".")
	// 	// if len(d.buffersArr) > 0 {
	// 	// 	bd := d.buffersArr[d.bufferToUpdate]
	// 	// 	d.bufferToUpdate = (d.bufferToUpdate + 1) % len(d.buffersArr)
	// 	// 	update := libevdiRequestUpdate(d.h, bd.buffer.id)
	// 	// 	if update {
	// 	// 		var numRects int32
	// 	// 		var rects *evdiRect
	// 	// 		libevdiGrabPixels(d.h, &rects, &numRects)
	// 	// 	}
	// 	// }
	// }
	return fmt.Errorf("TODO")
}

package libudev

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type LogFn func(udev UDev, priority int32, file string, line int32, fn string, format string, args unsafe.Pointer)

func SetLogFn(udev UDev, fn LogFn) {
	initLib()

	infoLck.Lock()
	defer infoLck.Unlock()

	track, ok := infoTrack[udev]
	if !ok {
		panic("invalid udev")
	}
	if track.logClosure != nil {
		ffi.ClosureFree(track.logClosure)
	}

	// allocate the closure function
	var callback unsafe.Pointer
	track.logClosure = ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)
	if track.logClosure == nil {
		panic("closure alloc failed")
	}

	// describe the closure's signature
	var cifCallback ffi.Cif
	/*!
	 * @brief Transfer a log message from libCEC to the client.
	 * @param cbparam             Callback parameter provided when the callbacks were set up
	 * @param message             The message to transfer.
	 */
	// void (*log_fn)(struct udev *udev, int priority, const char *file,
	//                int line, const char *fn, const char *format, va_list args))
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 7, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer,
		&ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fncb := ffi.NewCallback(logFnCallback)

	// prepare the closure
	userData := unsafe.Pointer(&fn)
	if status := ffi.PrepClosureLoc(track.logClosure, &cifCallback, fncb, userData, callback); status != ffi.OK {
		panic(status)
	}

	libudevSetLogFn(udev, callback)
}

func logFnCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {

	argArr := unsafe.Slice(args, cif.NArgs)

	udevPtr := *(*unsafe.Pointer)(argArr[0])

	udev := UDev(udevPtr)

	fn := *(*LogFn)(userData)

	if fn != nil {
		fn(udev, 0, "", 0, "", "", nil)
	}

	return 0
}

func SetLogPriority(udev UDev, priority LogPriority) {
	initLib()

	libudevSetLogPriority(udev, int32(priority))
}

func GetLogPriority(udev UDev) LogPriority {
	initLib()

	p := libudevGetLogPriority(udev)
	return LogPriority(p)
}

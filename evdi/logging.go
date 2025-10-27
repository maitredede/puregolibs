package evdi

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

var (
	loggingLck     sync.Mutex
	loggingClosure *ffi.Closure
)

//	struct evdi_logging {
//		void (*function)(void *user_data, const char *fmt, ...);
//		void *user_data;
//	};
type evdiLogging struct {
	function unsafe.Pointer
	userData unsafe.Pointer
}

func SetLogging(logger func(string)) {
	initLib()

	loggingLck.Lock()
	defer loggingLck.Unlock()

	if loggingClosure != nil {
		ffi.ClosureFree(loggingClosure)
		loggingClosure = nil
	}

	if logger == nil {
		return
	}

	var logCallback unsafe.Pointer
	loggingClosure = ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &logCallback)
	if loggingClosure == nil {
		panic("closure alloc failed")
	}
	// describe the closure's signature
	var cifCallback ffi.Cif
	// void (*function)(void *user_data, const char *fmt, ...);
	if status := ffi.PrepCifVar(&cifCallback, ffi.DefaultAbi, 2, 2, &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	// fn will be called, then the closure gets invoked
	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		argArr := unsafe.Slice(args, cif.NArgs)
		// userData := *(*unsafe.Pointer)(argArr[0])
		msgFormatPtr := *(**byte)(argArr[1])
		msgFormat := unix.BytePtrToString(msgFormatPtr)
		msgArgs := make([]any, 0)

		// TODO find a way for vararg callbacks

		msg := fmt.Sprintf(msgFormat, msgArgs...)

		logger(msg)
		return 0
	})

	// prepare the closure
	if status := ffi.PrepClosureLoc(loggingClosure, &cifCallback, fn, nil, logCallback); status != ffi.OK {
		panic(status)
	}

	log := evdiLogging{
		function: logCallback,
		userData: nil,
	}
	logType := ffi.NewType(&ffi.TypePointer, &ffi.TypePointer)

	setLoggingFn := mustGetSymbol("evdi_set_logging")
	var cifSetLogging ffi.Cif
	if status := ffi.PrepCif(&cifSetLogging, ffi.DefaultAbi, 1, &ffi.TypeVoid, &logType); status != ffi.OK {
		panic(status)
	}

	ffi.Call(&cifSetLogging, setLoggingFn, nil, unsafe.Pointer(&log))
}

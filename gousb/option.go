package gousb

import (
	"unsafe"
)

type libusbOption int32

const (
	optionLogLevel          libusbOption = 0
	optionUseUSBDK          libusbOption = 1
	optionNoDeviceDiscovery libusbOption = 2
	optionLogCB             libusbOption = 3
	optionMax               libusbOption = 4
)

const optionWeakAuthority libusbOption = optionNoDeviceDiscovery

type NativeLibusbInitOption struct {
	option libusbOption
	_      [4]byte
	val    [8]byte
}

type InitOption interface {
	build() NativeLibusbInitOption
}

type optLogLevel struct {
	level LogLevel
}

func (o *optLogLevel) build() NativeLibusbInitOption {
	no := NativeLibusbInitOption{
		option: optionLogLevel,
	}
	*(*int32)(unsafe.Pointer(&no.val[0])) = int32(o.level)
	return no
}

func WithLogLevel(level LogLevel) InitOption {
	return &optLogLevel{
		level: level,
	}
}

// type optLogCallback struct {
// 	cb LogCallback
// }

// func (o *optLogCallback) build() NativeLibusbInitOption {

// 	no := NativeLibusbInitOption{
// 		option: optionLogLevel,
// 	}

// 	nativeCallback := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
// 		argArr := unsafe.Slice(args, cif.NArgs)

// 		argLevelPtr := argArr[1]
// 		argStrPtr := argArr[2]

// 		level := *(*LogLevel)(argLevelPtr)
// 		str := "TODO"
// 		_ = argStrPtr

// 		o.cb(nil, level, str)

// 		return 0
// 	})

// 	*(*unsafe.Pointer)(unsafe.Pointer(&no.val[0])) = unsafe.Pointer(&nativeCallback)
// 	return no
// }

// func WithLogCallback(cb LogCallback) InitOption {
// 	return &optLogCallback{
// 		cb: cb,
// 	}
// }

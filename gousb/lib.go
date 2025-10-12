package gousb

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libusb.dylib"
	case "linux":
		return "libusb-1.0.so"
	case "windows":
		return "libusb-1.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libusbInit, initPtr, "libusb_init")
	// TODO : purego.RegisterLibFunc(&libusbInitContext, initPtr, "libusb_init_context")
	purego.RegisterLibFunc(&libusbExit, initPtr, "libusb_exit")
	purego.RegisterLibFunc(&libusbSetDebug, initPtr, "libusb_set_debug")
	// libusb_set_log_cb
	// libusb_get_version
	// libusb_has_capability
	purego.RegisterLibFunc(&libusbErrorName, initPtr, "libusb_error_name")
	purego.RegisterLibFunc(&libusbSetLocale, initPtr, "libusb_setlocale")
	purego.RegisterLibFunc(&libusbStrError, initPtr, "libusb_strerror")
}

var (
	libusbInit func(ctx *unsafe.Pointer) int32
	// TODO : libusbInitContext func(ctx *unsafe.Pointer) int32
	libusbExit     func(ctx unsafe.Pointer)
	libusbSetDebug func(ctx unsafe.Pointer, level LogLevel)
	// libusb_set_log_cb
	// libusb_get_version
	// libusb_has_capability
	libusbErrorName func(errorCode int32) string
	// libusbSetLocale func(ctx unsafe.Pointer, locale unsafe.Pointer) int32
	libusbSetLocale func(ctx unsafe.Pointer, locale string) int32
	libusbStrError  func(errorCode int32) string
)

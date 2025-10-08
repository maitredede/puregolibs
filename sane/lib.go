package sane

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/maitredede/puregolibs/sane/internal"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libfreefare.dylib"
	case "linux":
		return "libsane.so"
	case "windows":
		return "sane.dll"
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

var (
	libSaneInit func(versionCode *internal.SANE_Int, authorize uintptr) SANE_Status
	libSaneExit func()
	// libSaneGetDevices          func(deviceList *uintptr, localOnly SANE_Bool) SANE_Status
	libSaneOpen                func(name string, handle *internal.SANE_Handle) SANE_Status
	libSaneClose               func(h internal.SANE_Handle)
	libSaneGetOptionDescriptor func(h internal.SANE_Handle, n internal.SANE_Int) *internal.SANE_Option_Descriptor
	libSaneControlOption       func(h internal.SANE_Handle, n internal.SANE_Int, a SANE_Action, v unsafe.Pointer, i *internal.SANE_Int) SANE_Status
	libSaneGetParameters       func(h internal.SANE_Handle, p uintptr /*SANE_Parameters **/) SANE_Status
	libSaneStart               func(h internal.SANE_Handle) SANE_Status
	libSaneRead                func(h internal.SANE_Handle, buf *internal.SANE_Byte, maxLen internal.SANE_Int, ln *internal.SANE_Int) SANE_Status
	libSaneCancel              func(h internal.SANE_Handle)
	libSaneSetIOMode           func(h internal.SANE_Handle, m internal.SANE_Bool) SANE_Status
	libSaneGetSelectFD         func(h internal.SANE_Handle, fd *internal.SANE_Int) SANE_Status
	libSaneStrStatus           func(status SANE_Status) string
)

func libInitFuncs() {
	purego.RegisterLibFunc(&libSaneInit, initPtr, "sane_init")
	purego.RegisterLibFunc(&libSaneExit, initPtr, "sane_exit")
	// purego.RegisterLibFunc(&libSaneGetDevices, initPtr, "sane_get_devices")
	purego.RegisterLibFunc(&libSaneOpen, initPtr, "sane_open")
	purego.RegisterLibFunc(&libSaneClose, initPtr, "sane_close")
	purego.RegisterLibFunc(&libSaneGetOptionDescriptor, initPtr, "sane_get_option_descriptor")
	purego.RegisterLibFunc(&libSaneControlOption, initPtr, "sane_control_option")
	purego.RegisterLibFunc(&libSaneGetParameters, initPtr, "sane_get_parameters")
	purego.RegisterLibFunc(&libSaneStart, initPtr, "sane_start")
	purego.RegisterLibFunc(&libSaneRead, initPtr, "sane_read")
	purego.RegisterLibFunc(&libSaneCancel, initPtr, "sane_cancel")
	purego.RegisterLibFunc(&libSaneSetIOMode, initPtr, "sane_set_io_mode")
	purego.RegisterLibFunc(&libSaneGetSelectFD, initPtr, "sane_get_select_fd")
	purego.RegisterLibFunc(&libSaneStrStatus, initPtr, "sane_strstatus")
}

func Init() error {
	libInit()

	var versionCode *internal.SANE_Int = nil
	// authentication not supported yet
	var authorize uintptr = 0
	s := libSaneInit(versionCode, authorize)
	if s != StatusGood {
		return mkError(s)
	}
	return nil
}

func Exit() {
	libInit()

	libSaneExit()
}

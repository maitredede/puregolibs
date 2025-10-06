package sane

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

type SANE_Word = int32
type SANE_Int = SANE_Word
type SANE_Bool int32

const (
	SANE_FALSE SANE_Bool = 0
	SANE_TRUE  SANE_Bool = 1
)

func (b SANE_Bool) Go() bool {
	return b != SANE_FALSE
}

type SANE_Handle = uintptr
type SANE_Byte = byte

var (
	libSaneInit func(versionCode *SANE_Int, authorize uintptr) SANE_Status
	libSaneExit func()
	// libSaneGetDevices          func(deviceList *uintptr, localOnly SANE_Bool) SANE_Status
	libSaneOpen                func(name string, handle *SANE_Handle) SANE_Status
	libSaneClose               func(h SANE_Handle)
	libSaneGetOptionDescriptor func(h SANE_Handle, n SANE_Int) uintptr // const SANE_Option_Descriptor *
	libSaneControlOption       func(h SANE_Handle, n SANE_Int, a SANE_Action, v unsafe.Pointer, i *SANE_Int) SANE_Status
	libSaneGetParameters       func(h SANE_Handle, p uintptr /*SANE_Parameters **/) SANE_Status
	libSaneStart               func(h SANE_Handle) SANE_Status
	libSaneRead                func(h SANE_Handle, buf *SANE_Byte, maxLen SANE_Int, ln *SANE_Int) SANE_Status
	libSaneCancel              func(h SANE_Handle)
	libSaneSetIOMode           func(h SANE_Handle, m SANE_Bool) SANE_Status
	libSaneGetSelectFD         func(h SANE_Handle, fd *SANE_Int) SANE_Status
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

	var versionCode *SANE_Int = nil
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

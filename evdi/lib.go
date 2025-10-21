package evdi

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	libLck sync.Mutex
	libPtr uintptr
	libErr error
)

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(libPtr, sym)
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func initLib() {
	libLck.Lock()
	defer libLck.Unlock()

	if libErr != nil {
		panic(libErr)
	}
	if libPtr != 0 {
		return
	}

	libPtr, libErr = purego.Dlopen("libevdi.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libErr != nil {
		panic(libErr)
	}

	purego.RegisterLibFunc(&libEvdiGetLibVersion, libPtr, "evdi_get_lib_version")
	purego.RegisterLibFunc(&libEvdiCheckDevice, libPtr, "evdi_check_device")
	purego.RegisterLibFunc(&libEvdiAddDevice, libPtr, "evdi_add_device")
	purego.RegisterLibFunc(&libEvdiOpen, libPtr, "evdi_open")
	purego.RegisterLibFunc(&libEvdiClose, libPtr, "evdi_close")
	purego.RegisterLibFunc(&libXorgRunning, libPtr, "Xorg_running")
	// purego.RegisterLibFunc(&libEvdiSetLogging, libPtr, "evdi_set_logging")
}

var (
	libEvdiGetLibVersion func(version unsafe.Pointer)
	libEvdiCheckDevice   func(device int32) DeviceStatus
	libEvdiAddDevice     func() int32
	libEvdiOpen          func(device int32) handle
	libEvdiClose         func(handle handle) handle
	libXorgRunning       func() bool
	// libEvdiSetLogging       func(evdi) bool
)

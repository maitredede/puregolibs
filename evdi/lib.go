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
	purego.RegisterLibFunc(&libEvdiOpenAttachedToFixed, libPtr, "evdi_open_attached_to_fixed")
	purego.RegisterLibFunc(&libEvdiClose, libPtr, "evdi_close")
	purego.RegisterLibFunc(&libXorgRunning, libPtr, "Xorg_running")
	// purego.RegisterLibFunc(&libEvdiSetLogging, libPtr, "evdi_set_logging")

	purego.RegisterLibFunc(&libEvdiConnect, libPtr, "evdi_connect")
	purego.RegisterLibFunc(&libEvdiConnect2, libPtr, "evdi_connect2")
	purego.RegisterLibFunc(&libEvdiDisconnect, libPtr, "evdi_disconnect")
}

var (
	libEvdiGetLibVersion       func(version unsafe.Pointer)
	libEvdiCheckDevice         func(device int32) DeviceStatus
	libEvdiAddDevice           func() int32
	libEvdiOpen                func(device int32) evdiHandle
	libEvdiOpenAttachedToFixed func(sysfsParentDevice unsafe.Pointer, length uint32) evdiHandle
	libEvdiClose               func(handle evdiHandle)
	libXorgRunning             func() bool
	// libEvdiSetLogging       func(evdi) bool

	libEvdiConnect    func(handle evdiHandle, edid *byte, edidLength uint32, skuAreaLimit uint32)
	libEvdiConnect2   func(handle evdiHandle, edid *byte, edidLength uint32, pixelArealimit uint32, pixelPerSecondLimit uint32)
	libEvdiDisconnect func(handle evdiHandle)
)

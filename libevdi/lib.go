package libevdi

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if initPtr == 0 {
		name := "libevdi.so.1"

		initPtr, initError = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if initError != nil {
			err := fmt.Errorf("error loading library %s: %w", name, initError)
			panic(err)
		}

		libInitFuncs()
	}
}

func getSymbol(sym string) (uintptr, error) {
	return purego.Dlsym(initPtr, sym)
}

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&evdiOpen, initPtr, "evdi_open")
	purego.RegisterLibFunc(&evdiOpenAttachedToFixed, initPtr, "evdi_open_attached_to_fixed")
	purego.RegisterLibFunc(&evdiClose, initPtr, "evdi_close")

	purego.RegisterLibFunc(&evdiConnect, initPtr, "evdi_connect")
	purego.RegisterLibFunc(&evdiConnect2, initPtr, "evdi_connect2")
	purego.RegisterLibFunc(&evdiDisconnect, initPtr, "evdi_disconnect")
}

var (
	evdiOpen                func(device int) Handle
	evdiOpenAttachedToFixed func(sysfsParentDevice unsafe.Pointer, length int) Handle
	evdiClose               func(handle Handle)

	evdiConnect    func(handle Handle, edid unsafe.Pointer, edidLen uint, skuAreaLimit uint32)
	evdiConnect2   func(handle Handle, edid unsafe.Pointer, edidLen uint, pixelAreaLimit, pixelPerSecondLimit uint32)
	evdiDisconnect func(handle Handle)
)

package fontconfig

import (
	"fmt"
	"runtime"
	"sync"

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
		return "libfontconfig.so.1"
	case "windows":
		return "libfontconfig-1.dll"
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
	purego.RegisterLibFunc(&libFcGetVersion, initPtr, "FcGetVersion")

	purego.RegisterLibFunc(&libFcConfigHome, initPtr, "FcConfigHome")
	purego.RegisterLibFunc(&libFcConfigGetCurrent, initPtr, "FcConfigGetCurrent")
}

var (
	libFcGetVersion func() int32

	libFcConfigHome       func() string
	libFcConfigGetCurrent func() fcConfigPtr
)

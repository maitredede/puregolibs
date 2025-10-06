package cec

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
	// 	return "libfreefare.dylib"
	case "linux":
		return "libcec.so"
	case "windows":
		return "cec.dll"
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
	libCecInitialise func(configuration *NativeConfiguration) uintptr
	libCecDestroy    func(connection uintptr)

	libCecClearConfiguration      func(configuration *NativeConfiguration)
	libCecGetCurrentConfiguration func(connection uintptr, configuration *NativeConfiguration) int32

	libCecGetLibInfo func(connection uintptr) string
)

func libInitFuncs() {
	purego.RegisterLibFunc(&libCecInitialise, initPtr, "libcec_initialise")
	purego.RegisterLibFunc(&libCecDestroy, initPtr, "libcec_destroy")

	purego.RegisterLibFunc(&libCecClearConfiguration, initPtr, "libcec_clear_configuration")

	purego.RegisterLibFunc(&libCecGetCurrentConfiguration, initPtr, "libcec_get_current_configuration")

	purego.RegisterLibFunc(&libCecGetLibInfo, initPtr, "libcec_get_lib_info")
}

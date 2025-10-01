package sane

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

type SANE_Int = int16

var (
	libSaneInit func(versionCode *SANE_Int, authorize uintptr) SANE_Status
	libSaneExit func()
)

func libInitFuncs() {
	purego.RegisterLibFunc(&libSaneInit, initPtr, "sane_init")
	purego.RegisterLibFunc(&libSaneExit, initPtr, "sane_exit")
	purego.RegisterLibFunc(&libSaneStrStatus, initPtr, "sane_strstatus")
}

func Init() error {
	libInit()

	// authentication not supported yet
	s := libSaneInit(nil, 0)
	if s != StatusGood {
		return mkError(s)
	}
	return nil
}

func Exit() {
	libInit()

	libSaneExit()
}

package imagick

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce  sync.Mutex
	libWand      uintptr
	libWandError error
	libCore      uintptr
	libCoreError error
)

func getCoreLibrary() (string, error) {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libMagickCore.dylib"
	case "linux":
		return "libMagickCore-7.Q16.so", nil
	// case "windows":
	// 	return "libMagickCore-7.dll"
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}

func getWandLibrary() (string, error) {
	switch runtime.GOOS {
	// case "darwin":
	// 	return "libMagickWand.dylib"
	case "linux":
		return "libMagickWand-7.Q16.so", nil
	// case "windows":
	// 	return "libMagickWand-7.dll"
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if err := errors.Join(libCoreError, libWandError); err != nil {
		panic(err)
	}

	var initFuncs bool
	if libCore == 0 {
		var name string
		name, libCoreError = getCoreLibrary()
		if libCoreError == nil {
			libCore, libCoreError = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
			if libCoreError == nil {
				initFuncs = true
			}
		}
	}
	if libWand == 0 {
		var name string
		name, libWandError = getWandLibrary()
		if libWandError == nil {
			libWand, libWandError = purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		}
		if libWandError == nil {
			initFuncs = true
		}
	}
	if err := errors.Join(libCoreError, libWandError); err != nil {
		panic(err)
	}
	if !initFuncs {
		return
	}

	purego.RegisterLibFunc(&libCoreGetVersion, libCore, "GetMagickVersion")
	purego.RegisterLibFunc(&libCoreGenesis, libCore, "MagickCoreGenesis")
	purego.RegisterLibFunc(&libCoreTerminus, libCore, "MagickCoreTerminus")

	purego.RegisterLibFunc(&libWandGetVersion, libWand, "GetMagickVersion")
	purego.RegisterLibFunc(&libWandIsMagickCoreInstantiated, libWand, "IsMagickCoreInstantiated")
	purego.RegisterLibFunc(&libWandIsMagickWandInstantiated, libWand, "IsMagickWandInstantiated")
	purego.RegisterLibFunc(&libWandGenesis, libWand, "MagickWandGenesis")
	purego.RegisterLibFunc(&libWandTerminus, libWand, "MagickWandTerminus")
	purego.RegisterLibFunc(&libWandNewMagickWand, libWand, "NewMagickWand")
	purego.RegisterLibFunc(&libWandIsMagickWand, libWand, "IsMagickWand")
	purego.RegisterLibFunc(&libWandDestroyMagickWand, libWand, "DestroyMagickWand")
	purego.RegisterLibFunc(&libWandMagickReadImageBlob, libWand, "MagickReadImageBlob")
	purego.RegisterLibFunc(&libWandMagickSetImageFormat, libWand, "MagickSetImageFormat")
	purego.RegisterLibFunc(&libWandMagickGetImageBlob, libWand, "MagickGetImageBlob")
	purego.RegisterLibFunc(&libWandMagickRelinquishMemory, libWand, "MagickRelinquishMemory")
	purego.RegisterLibFunc(&libWandMagickGetInterlaceScheme, libWand, "MagickGetInterlaceScheme")
	purego.RegisterLibFunc(&libWandMagickSetInterlaceScheme, libWand, "MagickSetInterlaceScheme")
	purego.RegisterLibFunc(&libWandMagickGetException, libWand, "MagickGetException")
}

var (
	libCoreGetVersion func() string
	libCoreGenesis    func()
	libCoreTerminus   func()

	libWandGetVersion               func() string
	libWandIsMagickCoreInstantiated func() bool
	libWandIsMagickWandInstantiated func() bool
	libWandGenesis                  func()
	libWandTerminus                 func()
	libWandNewMagickWand            func() magickWandPtr
	libWandIsMagickWand             func(wand magickWandPtr) bool
	libWandDestroyMagickWand        func(wand magickWandPtr) magickWandPtr
	libWandMagickReadImageBlob      func(wand magickWandPtr, blob *byte, length uint32) bool
	libWandMagickSetImageFormat     func(wand magickWandPtr, format string) bool
	libWandMagickGetImageBlob       func(wand magickWandPtr, length *uint32) *byte
	libWandMagickRelinquishMemory   func(resource unsafe.Pointer)
	libWandMagickGetInterlaceScheme func(wand magickWandPtr) InterlaceType
	libWandMagickSetInterlaceScheme func(wand magickWandPtr, scheme InterlaceType) bool
	libWandMagickGetException       func(wand magickWandPtr, exceptionType *ExceptionType) string
)

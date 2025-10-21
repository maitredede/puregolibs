package imagick

import (
	"errors"
	"sync"
	"unsafe"
)

var (
	initLckOnce  sync.Mutex
	libWandError error
	libCoreError error
)

func libInit() {
	initLckOnce.Lock()
	defer initLckOnce.Unlock()

	if err := errors.Join(libCoreError, libWandError); err != nil {
		panic(err)
	}

	var initFuncs bool
	if !isLibLoaded(libCore) {
		libCore, libCoreError = loadLib(libCoreName)
		if libCoreError == nil {
			initFuncs = true
		}
	}
	if !isLibLoaded(libWand) {
		libWand, libWandError = loadLib(libWandName)
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

	RegisterLibFunc(&libCoreGetVersion, libCore, "GetMagickVersion")
	RegisterLibFunc(&libCoreGenesis, libCore, "MagickCoreGenesis")
	RegisterLibFunc(&libCoreTerminus, libCore, "MagickCoreTerminus")

	RegisterLibFunc(&libWandGetVersion, libWand, "GetMagickVersion")
	RegisterLibFunc(&libWandIsMagickCoreInstantiated, libWand, "IsMagickCoreInstantiated")
	RegisterLibFunc(&libWandIsMagickWandInstantiated, libWand, "IsMagickWandInstantiated")
	RegisterLibFunc(&libWandGenesis, libWand, "MagickWandGenesis")
	RegisterLibFunc(&libWandTerminus, libWand, "MagickWandTerminus")
	RegisterLibFunc(&libWandNewMagickWand, libWand, "NewMagickWand")
	RegisterLibFunc(&libWandIsMagickWand, libWand, "IsMagickWand")
	RegisterLibFunc(&libWandDestroyMagickWand, libWand, "DestroyMagickWand")
	RegisterLibFunc(&libWandMagickReadImageBlob, libWand, "MagickReadImageBlob")
	RegisterLibFunc(&libWandMagickSetImageFormat, libWand, "MagickSetImageFormat")
	RegisterLibFunc(&libWandMagickGetImageBlob, libWand, "MagickGetImageBlob")
	RegisterLibFunc(&libWandMagickRelinquishMemory, libWand, "MagickRelinquishMemory")
	RegisterLibFunc(&libWandMagickGetInterlaceScheme, libWand, "MagickGetInterlaceScheme")
	RegisterLibFunc(&libWandMagickSetInterlaceScheme, libWand, "MagickSetInterlaceScheme")
	RegisterLibFunc(&libWandMagickGetException, libWand, "MagickGetException")
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

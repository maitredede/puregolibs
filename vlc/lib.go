package vlc

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libvlcGetVersion, initPtr, "libvlc_get_version")
	purego.RegisterLibFunc(&libvlcGetCompiler, initPtr, "libvlc_get_compiler")
	purego.RegisterLibFunc(&libvlcGetChangeset, initPtr, "libvlc_get_changeset")

	purego.RegisterLibFunc(&libvlcErrmsg, initPtr, "libvlc_errmsg")

	purego.RegisterLibFunc(&libvlcNew, initPtr, "libvlc_new")
	purego.RegisterLibFunc(&libvlcRelease, initPtr, "libvlc_release")
	purego.RegisterLibFunc(&libvlcRetain, initPtr, "libvlc_retain")

}

var (
	libvlcGetVersion   func() string
	libvlcGetCompiler  func() string
	libvlcGetChangeset func() string

	libvlcErrmsg func() string

	libvlcNew     func(argc int32, argv unsafe.Pointer) libvlcInstance
	libvlcRelease func(instance libvlcInstance)
	libvlcRetain  func(instance libvlcInstance) libvlcInstance
)

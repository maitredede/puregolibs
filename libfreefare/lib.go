package libfreefare

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
		return "libfreefare.so"
	case "windows":
		return "libfreefare.dll"
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
	libFree func(ptr uintptr)
)

func libInitFuncs() {
	if _, err := getSymbol("freefare_version"); err != nil {
		libVersion = func() string {
			return "unknown (<=0.4.0)"
		}
	} else {
		purego.RegisterLibFunc(&libVersion, initPtr, "freefare_version")
	}

	purego.RegisterLibFunc(&libFree, initPtr, "free")

	purego.RegisterLibFunc(&libGetTags, initPtr, "freefare_get_tags")
	purego.RegisterLibFunc(&libFreeTags, initPtr, "freefare_free_tags")

	purego.RegisterLibFunc(&libGetTagType, initPtr, "freefare_get_tag_type")
	purego.RegisterLibFunc(&libGetTagFirendlyName, initPtr, "freefare_get_tag_friendly_name")
	purego.RegisterLibFunc(&libGetTagUID, initPtr, "freefare_get_tag_uid")
	purego.RegisterLibFunc(&libStrError, initPtr, "freefare_strerror")

	purego.RegisterLibFunc(&libMifareClassicConnect, initPtr, "mifare_classic_connect")
	purego.RegisterLibFunc(&libMifareClassicDisconnect, initPtr, "mifare_classic_disconnect")

	purego.RegisterLibFunc(&libMadRead, initPtr, "mad_read")
	purego.RegisterLibFunc(&libMadGetVersion, initPtr, "mad_get_version")
	purego.RegisterLibFunc(&libMadSetVersion, initPtr, "mad_set_version")

	purego.RegisterLibFunc(&libMifareApplicationRead, initPtr, "mifare_application_read")
}

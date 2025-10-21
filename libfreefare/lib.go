package libfreefare

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

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

var (
	libVersion func() string
	libFree    func(ptr uintptr)

	// MifareTag	*freefare_get_tags (nfc_device *device)
	//libGetTags func(nfcDevice uintptr) *uintptr
	libGetTags func(nfcDevice unsafe.Pointer) *uintptr
	// void		 freefare_free_tags (MifareTag *tags)
	// libFreeTags func(tag *uintptr)
	libFreeTags func(tags *uintptr)
	// void		 freefare_free_tag(FreefareTag tag)
	libFreeTag func(tag uintptr)

	// enum mifare_tag_type freefare_get_tag_type (MifareTag tag)
	libGetTagType func(tag uintptr) TagType
	// const char	*freefare_get_tag_friendly_name (MifareTag tag)
	libGetTagFirendlyName func(tag uintptr) string
	// char		*freefare_get_tag_uid (MifareTag tag)
	libGetTagUID func(tag uintptr) uintptr

	// const char	*freefare_strerror (MifareTag tag)
	libStrError func(tag uintptr) string

	// int		 mifare_classic_connect (MifareTag tag)
	libMifareClassicConnect func(tag uintptr) int16
	// int		 mifare_classic_disconnect (MifareTag tag)
	libMifareClassicDisconnect func(tag uintptr) int16

	// Mad		 mad_read (MifareTag tag)
	libMadRead func(tag uintptr) uintptr
	// int		 mad_get_version (Mad mad)
	libMadGetVersion func(mad uintptr) int16
	// void		 mad_set_version (Mad mad, const uint8_t version)
	libMadSetVersion func(mad uintptr, version byte)

	// ssize_t		 mifare_application_read (MifareTag tag, Mad mad, const MadAid aid, void *buf, size_t nbytes, const MifareClassicKey key, const MifareClassicKeyType key_type);
	libMifareApplicationRead func(tag uintptr, mad uintptr, aid *byte, buf uintptr, nbytes uint32, key *byte, keyType MifareClassicKeyType) int32
)

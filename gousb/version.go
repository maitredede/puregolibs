package gousb

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type libusbVersion struct {
	/** Library major version. */
	major uint16

	/** Library minor version. */
	minor uint16

	/** Library micro version. */
	micro uint16

	/** Library nano version. */
	nano uint16

	/** Library release candidate suffix string, e.g. "-rc4". */
	// const char *rc;
	rc unsafe.Pointer

	/** For ABI compatibility only. */
	// const char *describe;
	describe unsafe.Pointer
}

type Version struct {
	Major    int
	Minor    int
	Micro    int
	Nano     int
	RC       string
	Describe string
}

func GetVersion() Version {
	libInit()

	nativeV := libusbGetVersion()
	v := Version{
		Major: int(nativeV.major),
		Minor: int(nativeV.minor),
		Micro: int(nativeV.micro),
		Nano:  int(nativeV.nano),
	}
	if nativeV.rc != nil {
		v.RC = strings.GoString((*byte)(nativeV.rc))
	}
	if nativeV.describe != nil {
		v.Describe = strings.GoString((*byte)(nativeV.describe))
	}
	return v
}

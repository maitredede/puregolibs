package evdi

import (
	"fmt"
	"unsafe"
)

type nativeVersion struct {
	major int32
	minor int32
	patch int32
}

func VersionString() string {
	initLib()

	var v nativeVersion
	libEvdiGetLibVersion(unsafe.Pointer(&v))
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

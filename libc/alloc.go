package libc

import "unsafe"

func CAlloc(nmemb int32, size int32) unsafe.Pointer {
	initCLib()

	return libcCalloc(nmemb, size)
}

func Free(p unsafe.Pointer) {
	initCLib()

	libcFree(p)
}

package libc

import "unsafe"

func CAlloc(nmemb uint, size uint) unsafe.Pointer {
	initCLib()

	return libcCalloc(nmemb, size)
}

func Free(p unsafe.Pointer) {
	initCLib()

	libcFree(p)
}

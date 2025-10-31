//go:build unix

package strings

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// GoString copies a null-terminated char* to a Go string.
func GoString(p *byte) string {
	return unix.BytePtrToString(p)
}

// CString converts a go string to *byte that can be passed to C code.
func CString(s string) *byte {
	p, err := unix.BytePtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

// CStringL converts a go string to *byte that can be passed to C code, with its length
func CStringL(s string) (*byte, int) {
	a, err := unix.ByteSliceFromString(s)
	if err != nil {
		panic(err)
	}
	return &a[0], len(a)
}

// GoStringN copies a null-terminated char* to a Go string.
func GoStringN(p *byte, maxLen int) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
		if n >= maxLen {
			break
		}
	}

	return string(unsafe.Slice(p, n))
}

package strings

import (
	"unsafe"
)

// hasSuffix tests whether the string s ends with suffix.
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// CString converts a go string to *byte that can be passed to C code.
func CString(name string) *byte {
	if hasSuffix(name, "\x00") {
		return &(*(*[]byte)(unsafe.Pointer(&name)))[0]
	}
	b := make([]byte, len(name)+1)
	copy(b, name)
	return &b[0]
}

// CStringL converts a go string to *byte that can be passed to C code, and also returns its length
func CStringL(name string) (*byte, int) {
	if hasSuffix(name, "\x00") {
		l := len(name) - 1
		ptr := &(*(*[]byte)(unsafe.Pointer(&name)))[0]
		return ptr, l
	}

	l := len(name)
	b := make([]byte, l+1)
	copy(b, name)
	return &b[0], l
}

// GoString copies a null-terminated char* to a Go string.
func GoString(c uintptr) string {
	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}

// GoStringN copies a null-terminated char* to a Go string, or if reaching max len.
func GoStringN(c uintptr, max int) string {
	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		if length >= max {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}

// // GoStringB copies a null-terminated char* to a Go string.
// func GoStringB(arr []byte) string {
// 	if len(arr) == 0 {
// 		return ""
// 	}
// 	sb := strings.Builder{}
// 	sb.Grow(len(arr))
// 	for _, b := range arr {
// 		if b == 0 {
// 			break
// 		}
// 		if err := sb.WriteByte(b); err != nil {
// 			panic("FIXME")
// 		}
// 	}
// 	return sb.String()
// }

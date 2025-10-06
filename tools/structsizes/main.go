//go:build cgo

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var ptr uintptr
	fmt.Printf("nfc: sizeof uintptr: %d\n", unsafe.Sizeof(ptr))

	dumpNfc()
	dumpSane()
	dumpCec()

}

//go:build linux && cgo

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var ptr uintptr
	fmt.Printf("all: sizeof uintptr: %d\n", unsafe.Sizeof(ptr))

	dumpNfc()
	dumpSane()
	dumpCec()
	dumpUsb()
	dumpPlutobook()
}

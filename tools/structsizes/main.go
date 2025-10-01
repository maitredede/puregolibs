//go:build cgo

package main

// #include <nfc/nfc.h>
// #include <freefare.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var ptr uintptr
	fmt.Printf("sizeof uintptr: %d\n", unsafe.Sizeof(ptr))
	var nti C.nfc_target_info
	fmt.Printf("sizeof nfc_target_info: %d\n", unsafe.Sizeof(nti))

	var ffTag C.MifareTag
	fmt.Printf("sizeof MifareTag: %d\n", unsafe.Sizeof(ffTag))
}

//go:build linux && cgo

package main

// #include <nfc/nfc.h>
// #include <freefare.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func dumpNfc() {
	var nti C.nfc_target_info
	fmt.Printf("nfc: sizeof nfc_target_info: %d\n", unsafe.Sizeof(nti))

	var ffTag C.MifareTag
	fmt.Printf("freefare: sizeof MifareTag: %d\n", unsafe.Sizeof(ffTag))
}

//go:build cgo

package main

// #include <nfc/nfc.h>
// #include <freefare.h>
// #include <stdlib.h>
// #include <sane/sane.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var ptr uintptr
	fmt.Printf("nfc: sizeof uintptr: %d\n", unsafe.Sizeof(ptr))
	var nti C.nfc_target_info
	fmt.Printf("nfc: sizeof nfc_target_info: %d\n", unsafe.Sizeof(nti))

	var ffTag C.MifareTag
	fmt.Printf("freefare: sizeof MifareTag: %d\n", unsafe.Sizeof(ffTag))

	var saneInt C.SANE_Int
	fmt.Printf("sane: sizeof SANE_Int: %d\n", unsafe.Sizeof(saneInt))
	var saneAction C.SANE_Action
	fmt.Printf("sane: sizeof SANE_Action: %d\n", unsafe.Sizeof(saneAction))
	var saneHandle C.SANE_Handle
	fmt.Printf("sane: sizeof SANE_Handle: %d\n", unsafe.Sizeof(saneHandle))
	var saneBool C.SANE_Bool
	fmt.Printf("sane: sizeof SANE_Bool: %d\n", unsafe.Sizeof(saneBool))
	var saneByte C.SANE_Byte
	fmt.Printf("sane: sizeof SANE_Byte: %d\n", unsafe.Sizeof(saneByte))
	var saneChar C.SANE_Char
	fmt.Printf("sane: sizeof SANE_Char: %d\n", unsafe.Sizeof(saneChar))
}

//go:build linux && cgo

package main

// #include <sane/sane.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func dumpSane() {
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

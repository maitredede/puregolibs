package tools

import (
	"fmt"
	"unsafe"
)

func DumpMemory(ptr unsafe.Pointer, size uintptr) {
	bin := unsafe.Slice((*byte)(ptr), size)

	for i := 0; i < int(size); i++ {
		fmt.Printf(" %02x", bin[i])
	}
	fmt.Println()
}

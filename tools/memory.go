package tools

import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

func DumpMemory(ptr unsafe.Pointer, size uintptr) {
	bin := unsafe.Slice((*byte)(ptr), size)
	fmt.Println(strings.ToHexString(bin))
}

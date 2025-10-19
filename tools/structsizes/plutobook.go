//go:build linux && cgo

package main

// #include <plutobook.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/plutobook"
	"github.com/maitredede/puregolibs/tools"
)

func dumpPlutobook() {
	var nPageSize C.plutobook_page_size_t = C.plutobook_page_size_t{width: C.float(plutobook.PageSizeA4.Width), height: C.float(plutobook.PageSizeA4.Height)}
	var goPageSize plutobook.PageSize = plutobook.PageSizeA4
	fmt.Printf("plutobook: sizeof plutobook_page_size_t: %d plutobook.PageSize: %d\n", unsafe.Sizeof(nPageSize), unsafe.Sizeof(goPageSize))
	tools.DumpMemory(unsafe.Pointer(&nPageSize), unsafe.Sizeof(nPageSize))
	tools.DumpMemory(unsafe.Pointer(&goPageSize), unsafe.Sizeof(goPageSize))
}

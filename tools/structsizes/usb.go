//go:build linux && cgo

package main

// #include <libusb-1.0/libusb.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/gousb"
)

func dumpUsb() {
	var initOption C.struct_libusb_init_option
	var goInitOption gousb.NativeLibusbInitOption
	fmt.Printf("sane: sizeof libusb_init_option=%d goInitOption=%d\n", unsafe.Sizeof(initOption), unsafe.Sizeof(goInitOption))
}

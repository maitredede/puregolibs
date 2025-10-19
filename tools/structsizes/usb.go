//go:build linux && cgo

package main

// #include <libusb-1.0/libusb.h>
// #include <sys/time.h>
// #include <sys/types.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/gousb"
)

func dumpUsb() {
	var initOption C.struct_libusb_init_option
	var goInitOption gousb.NativeLibusbInitOption
	fmt.Printf("usb: sizeof libusb_init_option=%d goInitOption=%d\n", unsafe.Sizeof(initOption), unsafe.Sizeof(goInitOption))

	var timeval C.struct_timeval
	var goTimeval gousb.NativeTimeval
	fmt.Printf("usb: sizeof timeval=%d goTimeval=%d\n", unsafe.Sizeof(timeval), unsafe.Sizeof(goTimeval))
}

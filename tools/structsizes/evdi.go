//go:build ignore

package main

// #include <evdi_lib.h>
// // #cgo LDFLAGS: -llibevdi
// // #cgo pkg-config: libevdi
import "C"
import "fmt"

func dumpEvdi() {
	var v C.struct_evdi_lib_version
	C.evdi_get_lib_version(&v)

	fmt.Printf("evdi: version=%d.%d.%d\n", v.version_major, v.version_minor, v.version_patchlevel)
}

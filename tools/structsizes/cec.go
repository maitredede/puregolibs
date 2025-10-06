//go:build cgo

package main

// #include <libcec/cecc.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/cec"
)

func dumpCec() {
	var cecConfig C.libcec_configuration
	var puregoConfig cec.NativeConfiguration
	fmt.Printf("cec: sizeof libcec_configuration: %d (purego:%d)\n", unsafe.Sizeof(cecConfig), unsafe.Sizeof(puregoConfig))
}

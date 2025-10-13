//go:build linux && cgo

package main

// #include <fontconfig/fontconfig.h>
import "C"

func dumpFontconfig() {
	C.FcConfigHome()
	C.FcGetVersion()
}

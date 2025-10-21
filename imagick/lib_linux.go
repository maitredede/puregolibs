//go:build linux

package imagick

import "github.com/ebitengine/purego"

const (
	libCoreName string = "libMagickCore-7.Q16.so"
	libWandName string = "libMagickWand-7.Q16.so"
)

var (
	libWand uintptr
	libCore uintptr
)

func isLibLoaded(lib uintptr) bool {
	return lib != 0
}

func loadLib(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

func RegisterLibFunc(fptr any, handle uintptr, name string) {
	purego.RegisterLibFunc(fptr, handle, name)
}

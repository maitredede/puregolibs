//go:build windows

package imagick

import (
	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

const (
	libCoreName string = "libMagickCore-7.Q16.dll"
	libWandName string = "libMagickWand-7.Q16.dll"
)

var (
	libWand *windows.DLL
	libCore *windows.DLL
)

func isLibLoaded(lib *windows.DLL) bool {
	return lib != nil
}

func loadLib(name string) (*windows.DLL, error) {
	return windows.LoadDLL(name)
}

func RegisterLibFunc(fptr any, dll *windows.DLL, name string) {
	sym, err := dll.FindProc(name)
	if err != nil {
		panic(err)
	}
	purego.RegisterFunc(fptr, sym.Addr())
}

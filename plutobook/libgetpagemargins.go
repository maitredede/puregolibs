package plutobook

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	libGetPageMargins func(book uintptr) PageMargins
)

func registerFFIGetPageMargins() {
	libGetPageMarginsSym := mustGetSymbol("plutobook_get_page_margins")

	argTypes := []*ffi.Type{
		&ffi.TypePointer,
	}
	var libGetPageMarginsCIF ffi.Cif
	if ok := ffi.PrepCif(&libGetPageMarginsCIF, ffi.DefaultAbi, uint32(len(argTypes)), &ffiPageMarginsType, argTypes...); ok != ffi.OK {
		panic("plutobook_get_page_margins cif prep is not OK")
	}

	libGetPageMargins = func(book uintptr) PageMargins {
		var ret PageMargins
		args := []unsafe.Pointer{
			unsafe.Pointer(book),
		}
		ffi.Call(&libGetPageMarginsCIF, libGetPageMarginsSym, unsafe.Pointer(&ret), args...)
		return ret
	}
}

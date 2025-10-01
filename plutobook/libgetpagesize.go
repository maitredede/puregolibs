package plutobook

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	libGetPageSize   func(book uintptr) PageSize
	libGetPageSizeAt func(book uintptr, index int) PageSize
)

func registerFFIGetPageSize() {
	libGetPageSizeSym := mustGetSymbol("plutobook_get_page_size")

	var libGetPageSizeCIF ffi.Cif
	if ok := ffi.PrepCif(&libGetPageSizeCIF, ffi.DefaultAbi, 1, &ffiPageSizeType, &ffi.TypePointer); ok != ffi.OK {
		panic("plutobook_get_page_size cif prep is not OK")
	}

	libGetPageSize = func(book uintptr) PageSize {
		var ret PageSize
		args := []unsafe.Pointer{
			unsafe.Pointer(book),
		}
		ffi.Call(&libGetPageSizeCIF, libGetPageSizeSym, unsafe.Pointer(&ret), args...)
		return ret
	}
}

func registerFFIGetPageSizeAt() {
	sym := mustGetSymbol("plutobook_get_page_size_at")

	var cif ffi.Cif
	if ok := ffi.PrepCif(&cif, ffi.DefaultAbi, 2, &ffiPageSizeType, &ffi.TypePointer, &ffi.TypeUint16); ok != ffi.OK {
		panic("plutobook_get_page_size_at cif prep is not OK")
	}

	libGetPageSizeAt = func(book uintptr, index int) PageSize {
		nIndex := uint16(index)
		var ret PageSize
		args := []unsafe.Pointer{
			unsafe.Pointer(book),
			unsafe.Pointer(&nIndex),
		}
		ffi.Call(&cif, sym, unsafe.Pointer(&ret), args...)
		return ret
	}
}

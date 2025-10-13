package plutobook

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

func registerFFICreate() {
	libCreateSym := mustGetSymbol("plutobook_create")

	argTypes := []*ffi.Type{
		&ffiPageSizeType,
		&ffiPageMarginsType,
		&ffiMediaTypeType,
	}
	var libCreateCIF ffi.Cif
	if ok := ffi.PrepCif(&libCreateCIF, ffi.DefaultAbi, uint32(len(argTypes)), &ffi.TypePointer, argTypes...); ok != ffi.OK {
		panic("plutobook_create cif prep is not OK")
	}

	libCreate = func(pageSize PageSize, margins PageMargins, mediaType MediaType) bookPtr {
		var ret bookPtr
		args := []unsafe.Pointer{
			unsafe.Pointer(&pageSize),
			unsafe.Pointer(&margins),
			unsafe.Pointer(&mediaType),
		}
		ffi.Call(&libCreateCIF, libCreateSym, unsafe.Pointer(&ret), args...)
		return ret
	}
}

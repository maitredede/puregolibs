package plutobook

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	libCreate func(pageSize PageSize, margins PageMargins, mediaType MediaType) uintptr
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

	libCreate = func(pageSize PageSize, margins PageMargins, mediaType MediaType) uintptr {
		var aMediaType int16 = int16(mediaType)
		var ret uintptr
		args := []unsafe.Pointer{
			unsafe.Pointer(&pageSize),
			unsafe.Pointer(&margins),
			unsafe.Pointer(&aMediaType),
		}
		ffi.Call(&libCreateCIF, libCreateSym, unsafe.Pointer(&ret), args...)
		return ret
	}
}

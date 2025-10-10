package plutobook

import "github.com/jupiterrider/ffi"

type MediaType int32

const (
	MediaTypePrint MediaType = iota
	MediaTypeScreen
)

var (
	ffiMediaTypeType = ffi.TypeSint32

	libGetMediaType func(book uintptr) int32
)

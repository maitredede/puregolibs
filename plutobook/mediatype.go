package plutobook

import "github.com/jupiterrider/ffi"

type MediaType int

const (
	MediaTypePrint MediaType = iota
	MediaTypeScreen
)

var (
	ffiMediaTypeType = ffi.TypeSint16

	libGetMediaType func(book uintptr) int16
)

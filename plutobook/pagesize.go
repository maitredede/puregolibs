package plutobook

import (
	"github.com/jupiterrider/ffi"
)

var (
	ffiPageSizeType = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat)
)

type PageSize struct {
	// _ structs.HostLayout

	Width  float32
	Height float32
}

var (
	PageSizeA4 = PageSize{
		Width:  210.0 * float32(UnitsMM),
		Height: 297.0 * float32(UnitsMM),
	}
)

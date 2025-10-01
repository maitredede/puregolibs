package plutobook

import (
	"github.com/jupiterrider/ffi"
)

var (
	ffiPageMarginsType = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat)
	// ffiPageMarginsType = ffi.NewType(&ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble, &ffi.TypeDouble)
)

type PageMargins struct {
	// _ structs.HostLayout

	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

var (
	PageMarginsNone     = PageMargins{Top: 0, Right: 0, Bottom: 0, Left: 0}
	PageMarginsNormal   = PageMargins{Top: 72, Right: 72, Bottom: 72, Left: 72}
	PageMarginsNarrow   = PageMargins{Top: 36, Right: 36, Bottom: 36, Left: 36}
	PageMarginsModerate = PageMargins{Top: 72, Right: 54, Bottom: 72, Left: 54}
	PageMarginsWide     = PageMargins{Top: 72, Right: 144, Bottom: 72, Left: 144}
)

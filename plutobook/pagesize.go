package plutobook

import "github.com/jupiterrider/ffi"

var (
	ffiPageSizeType = ffi.NewType(&ffi.TypeFloat, &ffi.TypeFloat)
)

type PageSize struct {
	Width  float32
	Height float32
}

var (
	PageSizeNone = PageSize{
		Width:  0,
		Height: 0,
	}
	PageSizeA3 = PageSize{
		Width:  297.0 * float32(UnitsMM),
		Height: 420.0 * float32(UnitsMM),
	}
	PageSizeA4 = PageSize{
		Width:  210.0 * float32(UnitsMM),
		Height: 297.0 * float32(UnitsMM),
	}
	PageSizeA5 = PageSize{
		Width:  148.0 * float32(UnitsMM),
		Height: 210.0 * float32(UnitsMM),
	}
	PageSizeB4 = PageSize{
		Width:  250.0 * float32(UnitsMM),
		Height: 353.0 * float32(UnitsMM),
	}
	PageSizeB5 = PageSize{
		Width:  176.0 * float32(UnitsMM),
		Height: 250.0 * float32(UnitsMM),
	}
	PageSizeLetter = PageSize{
		Width:  8.5 * float32(UnitsIN),
		Height: 11.0 * float32(UnitsIN),
	}
	PageSizeLegal = PageSize{
		Width:  8.5 * float32(UnitsIN),
		Height: 14.0 * float32(UnitsIN),
	}
	PageSizeLedger = PageSize{
		Width:  11.0 * float32(UnitsIN),
		Height: 17.0 * float32(UnitsIN),
	}
)

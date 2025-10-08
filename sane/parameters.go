package sane

import (
	"unsafe"

	"github.com/maitredede/puregolibs/sane/internal"
)

type Parameters struct {
	Format        Format
	IsLastFrame   bool
	BytesPerLine  int
	PixelsPerLine int
	Lines         int
	Depth         int
}

type internalParameters struct {
	format          Format
	last_frame      internal.SANE_Bool
	bytes_per_line  internal.SANE_Int
	pixels_per_line internal.SANE_Int
	lines           internal.SANE_Int
	depth           internal.SANE_Int
}

func (h *Handle) GetParameters() (Parameters, error) {
	var np internalParameters
	ret := libSaneGetParameters(h.h, uintptr(unsafe.Pointer(&np)))
	if ret != StatusGood {
		return Parameters{}, mkError(ret)
	}
	p := Parameters{
		Format:        np.format,
		IsLastFrame:   np.last_frame.Go(),
		BytesPerLine:  int(np.bytes_per_line),
		PixelsPerLine: int(np.pixels_per_line),
		Lines:         int(np.lines),
		Depth:         int(np.depth),
	}
	return p, nil
}

package sane

import "unsafe"

type Parameters struct {
	Format        Frame
	IsLastFrame   bool
	BytesPerLine  int
	PixelsPerLine int
	Lines         int
	Depth         int
}

type internalParameters struct {
	format          Frame
	last_frame      SANE_Bool
	bytes_per_line  SANE_Int
	pixels_per_line SANE_Int
	lines           SANE_Int
	depth           SANE_Int
}

func GetParameters(h SANE_Handle) (Parameters, error) {
	var np internalParameters
	ret := libSaneGetParameters(h, uintptr(unsafe.Pointer(&np)))
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

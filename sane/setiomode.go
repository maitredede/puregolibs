package sane

import "github.com/maitredede/puregolibs/sane/internal"

func (h *Handle) SetIOMode(nonBlockingMode bool) error {
	var m internal.SANE_Bool = internal.SANE_FALSE
	if nonBlockingMode {
		m = internal.SANE_TRUE
	}
	ret := libSaneSetIOMode(h.h, m)
	if ret != StatusGood {
		return mkError(ret)
	}
	return nil
}

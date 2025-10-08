package sane

import (
	"unsafe"

	"github.com/maitredede/puregolibs/sane/internal"
)

func (h *Handle) GetSelectFD() (uintptr, error) {

	var fd internal.SANE_Int
	ret := libSaneGetSelectFD(h.h, &fd)
	if ret != StatusGood {
		return 0, mkError(ret)
	}
	v := *(*uintptr)(unsafe.Pointer(&fd))
	return v, nil
}

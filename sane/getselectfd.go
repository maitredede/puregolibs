package sane

import "unsafe"

func GetSelectFD(h SANE_Handle) (uintptr, error) {

	var fd SANE_Int
	ret := libSaneGetSelectFD(h, &fd)
	if ret != StatusGood {
		return 0, mkError(ret)
	}
	v := *(*uintptr)(unsafe.Pointer(&fd))
	return v, nil
}

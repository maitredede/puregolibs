package sane

import "github.com/maitredede/puregolibs/sane/internal"

type Handle struct {
	h internal.SANE_Handle
}

func (h *Handle) Native() uintptr {
	return h.h
}

func Open(name string) (*Handle, error) {
	var handle internal.SANE_Handle
	ret := libSaneOpen(name, &handle)
	if ret != StatusGood {
		return nil, mkError(ret)
	}
	h := &Handle{
		h: handle,
	}
	return h, nil
}

func (h *Handle) Close() {
	if h.h == 0 {
		return
	}
	libSaneClose(h.h)
	h.h = 0
}

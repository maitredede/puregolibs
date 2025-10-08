package sane

func (h *Handle) Start() error {
	ret := libSaneStart(h.h)
	if ret != StatusGood {
		return mkError(ret)
	}
	return nil
}

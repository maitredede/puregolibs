package sane

func (h *Handle) Cancel() {
	libSaneCancel(h.h)
}

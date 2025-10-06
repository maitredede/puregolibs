package sane

func Start(h SANE_Handle) error {
	ret := libSaneStart(h)
	if ret != StatusGood {
		return mkError(ret)
	}
	return nil
}

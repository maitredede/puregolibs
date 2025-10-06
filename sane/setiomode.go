package sane

func SetIOMode(h SANE_Handle, nonBlockingMode bool) error {
	var m SANE_Bool = SANE_FALSE
	if nonBlockingMode {
		m = SANE_TRUE
	}
	ret := libSaneSetIOMode(h, m)
	if ret != StatusGood {
		return mkError(ret)
	}
	return nil
}

package sane

func Open(name string) (SANE_Handle, error) {
	// cName := strings.CString(name)

	var handle SANE_Handle
	ret := libSaneOpen(name, &handle)
	if ret != StatusGood {
		return 0, mkError(ret)
	}
	return handle, nil
}

func Close(h SANE_Handle) {
	libSaneClose(h)
}

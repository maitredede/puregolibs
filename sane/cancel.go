package sane

func Cancel(h SANE_Handle) {
	libSaneCancel(h)
}

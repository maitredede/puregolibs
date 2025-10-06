package sane

import "io"

type saneStatusError struct {
	s SANE_Status
}

func (e *saneStatusError) Error() string {
	return libSaneStrStatus(e.s)
}

func mkError(s SANE_Status) error {
	if s == StatusEOF {
		return io.EOF
	}
	return &saneStatusError{
		s: s,
	}
}

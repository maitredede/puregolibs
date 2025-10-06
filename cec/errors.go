package cec

var ErrConnectionIsClosed *ErrorConnectionIsClosed = &ErrorConnectionIsClosed{m: "connection is closed"}

type ErrorConnectionIsClosed struct {
	m string
}

var _ error = (*ErrorConnectionIsClosed)(nil)

func (e ErrorConnectionIsClosed) Error() string {
	return e.m
}

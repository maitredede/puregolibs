package libfreefare

var ErrTagClosed *ErrorTagClosed = &ErrorTagClosed{m: "tag is closed"}

type ErrorTagClosed struct {
	m string
}

var _ error = (*ErrorTagClosed)(nil)

func (e ErrorTagClosed) Error() string {
	return e.m
}

package plutobook

var ErrBookIsClosed *ErrorBookIsClosed = &ErrorBookIsClosed{m: "book is closed"}

type ErrorBookIsClosed struct {
	m string
}

var _ error = (*ErrorBookIsClosed)(nil)

func (e ErrorBookIsClosed) Error() string {
	return e.m
}

var ErrCanvasIsClosed *ErrorCanvasIsClosed = &ErrorCanvasIsClosed{m: "canvas is closed"}

type ErrorCanvasIsClosed struct {
	m string
}

var _ error = (*ErrorCanvasIsClosed)(nil)

func (e ErrorCanvasIsClosed) Error() string {
	return e.m
}

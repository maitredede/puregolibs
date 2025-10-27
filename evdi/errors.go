package evdi

var ErrDeviceIsClosed *ErrorDeviceIsClosed = &ErrorDeviceIsClosed{m: "device is closed"}

type ErrorDeviceIsClosed struct {
	m string
}

var _ error = (*ErrorDeviceIsClosed)(nil)

func (e ErrorDeviceIsClosed) Error() string {
	return e.m
}

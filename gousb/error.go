package gousb

type USBError struct {
	value int32
	name  string
	msg   string
}

func (e *USBError) Error() string {
	return e.msg
}

func (e *USBError) Name() string {
	return e.name
}

func (e *USBError) Value() int32 {
	return e.value
}

func errorFromRet(ret int32) error {
	if ret == 0 {
		return nil
	}

	e := &USBError{
		value: ret,
	}
	e.name = libusbErrorName(ret)
	e.msg = libusbStrError(ret)
	return e
}

type InvalidContextError struct {
	msg string
}

func (e *InvalidContextError) Error() string {
	return e.msg
}

var ErrInvalidContext error = &InvalidContextError{msg: "invalid context"}

type InvalidDeviceError struct {
	msg string
}

func (e *InvalidDeviceError) Error() string {
	return e.msg
}

var ErrInvalidDevice error = &InvalidDeviceError{msg: "invalid device"}

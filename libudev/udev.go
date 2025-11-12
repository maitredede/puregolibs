package libudev

import "unsafe"

type UDev unsafe.Pointer

func New() UDev {
	initLib()

	return libudevNew()
}

func Unref(udev UDev) {
	initLib()

	libudevUnref(udev)
}

package evdi

import (
	"fmt"
	"unsafe"
)

type handle unsafe.Pointer

func AddDevice() int {
	initLib()
	return int(libEvdiAddDevice())
}

type Device struct {
	num int
	h   handle
}

func OpenDevice(device int) (*Device, error) {
	initLib()

	h := libEvdiOpen(int32(device))
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		num: device,
		h:   h,
	}
	return d, nil
}

package evdi

import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type evdiHandle unsafe.Pointer

func AddDevice() int {
	initLib()
	return int(libEvdiAddDevice())
}

type Device struct {
	h evdiHandle
}

func OpenDevice(device int) (*Device, error) {
	initLib()

	h := libEvdiOpen(int32(device))
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		h: h,
	}
	return d, nil
}

func OpenAttachedTo(sysfsParent string) (*Device, error) {
	initLib()

	cParent, l := strings.CStringL(sysfsParent)
	h := libEvdiOpenAttachedToFixed(unsafe.Pointer(cParent), uint32(l))
	if h == nil {
		return nil, fmt.Errorf("failed to open device")
	}
	d := &Device{
		h: h,
	}
	return d, nil
}

package libevdev

import (
	"fmt"
	"unsafe"
)

type Evdev unsafe.Pointer

/*New Initialize a new Evdev device. This function only allocates the
 * required memory and initializes the struct to sane default values.
 * To actually hook up the device to a kernel device, use
 * SetFd().
 *
 * Memory allocated through New() must be released by the
 * caller with Free().
 */
func New() Evdev {
	initLib()

	return libevdevNew()
}

func Free(evdev Evdev) {
	initLib()

	libevdevFree(evdev)
}

func SetFd(evdev Evdev, fd uintptr) error {
	initLib()

	ret := libevdevSetFd(evdev, fd)
	if ret == 0 {
		return nil
	}
	return fmt.Errorf("TODO error %v", ret)
}

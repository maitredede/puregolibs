package libevdev

import (
	"sync"

	"github.com/ebitengine/purego"
)

var (
	libLck sync.Mutex
	libPtr uintptr
	libErr error
)

func Initialize() error {
	libLck.Lock()
	defer libLck.Unlock()

	initLibNoPanic()
	return libErr
}

func Unload() error {
	libLck.Lock()
	defer libLck.Unlock()

	if libPtr == 0 {
		return nil
	}
	if err := purego.Dlclose(libPtr); err != nil {
		return err
	}
	libPtr = 0
	libErr = nil
	return nil
}

func initLib() {
	libLck.Lock()
	defer libLck.Unlock()

	initLibNoPanic()
	if libErr != nil {
		panic(libErr)
	}
}

func initLibNoPanic() {
	if libErr != nil {
		return
	}
	if libPtr != 0 {
		return
	}

	libPtr, libErr = purego.Dlopen("libevdev.so.2", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libErr != nil {
		return
	}

	purego.RegisterLibFunc(&libevdevNew, libPtr, "libevdev_new")
	purego.RegisterLibFunc(&libevdevFree, libPtr, "libevdev_free")
	purego.RegisterLibFunc(&libevdevSetFd, libPtr, "libevdev_set_fd")

	purego.RegisterLibFunc(&libevdevGetDriverVersion, libPtr, "libevdev_get_driver_version")
	purego.RegisterLibFunc(&libevdevGetName, libPtr, "libevdev_get_name")
	purego.RegisterLibFunc(&libevdevGetIDBusType, libPtr, "libevdev_get_id_bustype")
	purego.RegisterLibFunc(&libevdevGetIDVendor, libPtr, "libevdev_get_id_vendor")
	purego.RegisterLibFunc(&libevdevGetIDProduct, libPtr, "libevdev_get_id_product")
	purego.RegisterLibFunc(&libevdevGetIDVersion, libPtr, "libevdev_get_id_version")

	purego.RegisterLibFunc(&libevdevGrab, libPtr, "libevdev_grab")
}

var (
	libevdevNew   func() Evdev
	libevdevFree  func(evdev Evdev)
	libevdevSetFd func(evdev Evdev, fd uintptr) uintptr

	libevdevGetDriverVersion func(evdev Evdev) int32
	libevdevGetName          func(evdev Evdev) string
	libevdevGetIDBusType     func(evdev Evdev) uint16
	libevdevGetIDVendor      func(evdev Evdev) uint16
	libevdevGetIDProduct     func(evdev Evdev) uint16
	libevdevGetIDVersion     func(evdev Evdev) uint16

	libevdevGrab func(evdev Evdev, grab bool) int32
)

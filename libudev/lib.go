package libudev

import (
	"sync"
	"unsafe"

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

	libPtr, libErr = purego.Dlopen("libudev.so.1", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if libErr != nil {
		return
	}

	purego.RegisterLibFunc(&libudevNew, libPtr, "udev_new")
	purego.RegisterLibFunc(&libudevRef, libPtr, "udev_ref")
	purego.RegisterLibFunc(&libudevUnref, libPtr, "udev_unref")
	purego.RegisterLibFunc(&libudevSetLogFn, libPtr, "udev_set_log_fn")
	purego.RegisterLibFunc(&libudevGetLogPriority, libPtr, "udev_get_log_priority")
	purego.RegisterLibFunc(&libudevSetLogPriority, libPtr, "udev_set_log_priority")
	purego.RegisterLibFunc(&libudevGetUserData, libPtr, "udev_get_userdata")
	purego.RegisterLibFunc(&libudevSetUserData, libPtr, "udev_set_userdata")
}

var (
	libudevNew            func() UDev
	libudevRef            func(udev UDev) UDev
	libudevUnref          func(udev UDev)
	libudevSetLogFn       func(udev UDev, fn unsafe.Pointer)
	libudevGetLogPriority func(udev UDev) int32
	libudevSetLogPriority func(udev UDev, priority int32)
	libudevGetUserData    func(udev UDev) unsafe.Pointer
	libudevSetUserData    func(udev UDev, userData unsafe.Pointer)
)

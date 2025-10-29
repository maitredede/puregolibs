package drm

import (
	"os"
	"syscall"

	"github.com/maitredede/puregolibs/evdi/libevdi/drm/ioctl"
)

func IsMaster(f *os.File) bool {
	res := AuthMagic(f, 0)
	return res != syscall.EACCES
}

func DropMaster(file *os.File) error {
	err := ioctl.Do(uintptr(file.Fd()), uintptr(IOCTLDropMaster), 0)
	return err
}

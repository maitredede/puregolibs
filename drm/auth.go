//go:build linux

package drm

import (
	"os"
	"unsafe"

	"github.com/maitredede/puregolibs/drm/ioctl"
)

type DrmMagic uint32

type drmAuth struct {
	magic DrmMagic
}

func AuthMagic(file *os.File, magic DrmMagic) error {
	var auth drmAuth

	auth.magic = magic

	err := ioctl.Do(uintptr(file.Fd()), uintptr(IOCTLAuthMagic), uintptr(unsafe.Pointer(&auth)))

	return err
}

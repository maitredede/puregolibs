package drm

import (
	"unsafe"

	"github.com/maitredede/puregolibs/evdi/libevdi/drm/ioctl"
)

const IOCTLBase = 'd'

var (
	// DRM_IOWR(0x00, struct drm_version)
	IOCTLVersion = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmVersion{})), IOCTLBase, 0x00)
	// DRM_IOWR(0x01, struct drm_unique)
	IOCTLGetUnique = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmUnique{})), IOCTLBase, 0x01)

	// DRM_IOWR(0x0c, struct drm_get_cap)
	IOCTLGetCap = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmGetCap{})), IOCTLBase, 0x0c)
)

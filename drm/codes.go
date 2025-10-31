//go:build linux

package drm

import (
	"unsafe"

	"github.com/maitredede/puregolibs/drm/ioctl"
)

const (
	IOCTLBase   = 'd'
	CommandBase = 0x40
	CommandEnd  = 0xA0
)

var (
	// DRM_IOWR(0x00, struct drm_version)
	IOCTLVersion = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmVersion{})), IOCTLBase, 0x00)
	// DRM_IOWR(0x01, struct drm_unique)
	IOCTLGetUnique = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmUnique{})), IOCTLBase, 0x01)
	// DRM_IOR( 0x02, struct drm_auth)
	IOCTLGetMagic = ioctl.NewCode(ioctl.Read, uint16(unsafe.Sizeof(drmAuth{})), IOCTLBase, 0x02)
	// DRM_IOWR(0x03, struct drm_irq_busid)
	IOCTLIrqBusID = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmIrqBusID{})), IOCTLBase, 0x03)

	// DRM_IOW( 0x11, struct drm_auth)
	IOCTLAuthMagic = ioctl.NewCode(ioctl.Write, uint16(unsafe.Sizeof(drmAuth{})), IOCTLBase, 0x11)
	// IOCTLAuthMagic = drmIOW(0x11, uint16(unsafe.Sizeof(drmAuth{})))

	// DRM_IOWR(0x0c, struct drm_get_cap)
	IOCTLGetCap = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(drmGetCap{})), IOCTLBase, 0x0c)

	// DRM_IO(0x1e)
	IOCTLSetMaster = ioctl.NewCode(ioctl.None, 0, IOCTLBase, 0x1e)
	// DRM_IO(0x1f)
	IOCTLDropMaster = ioctl.NewCode(ioctl.None, 0, IOCTLBase, 0x1f)

	// DRM_IOCTL_MODE_MAP_DUMB    DRM_IOWR(0xB3, struct drm_mode_map_dumb)
	IOCTLModeMapDumb = ioctl.NewCode(ioctl.Read|ioctl.Write, uint16(unsafe.Sizeof(ModeMapDumb{})), IOCTLBase, 0xb3)
)

func drmIOW(nr, typ uint16) uint32 {
	return ioctl.IOW(IOCTLBase, nr, typ)
}

func drmIOR(nr, typ uint16) uint32 {
	return ioctl.IOR(IOCTLBase, nr, typ)
}

func drmIOWR(nr, typ uint16) uint32 {
	return ioctl.IOWR(IOCTLBase, nr, typ)
}

func drmIO(nr uint16) uint32 {
	return ioctl.IO(IOCTLBase, nr)
}

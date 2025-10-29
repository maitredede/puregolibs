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
)

func drmIOW(nr, typ uint16) uint32 {
	return ioctl.IOW(IOCTLBase, nr, typ)
}

func drmIOR(nr, typ uint16) uint32 {
	return ioctl.IOR(IOCTLBase, nr, typ)
}

func drmIORW(nr, typ uint16) uint32 {
	return ioctl.IORW(IOCTLBase, nr, typ)
}

func drmIO(nr uint16) uint32 {
	return ioctl.IO(IOCTLBase, nr)
}

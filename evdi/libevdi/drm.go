package libevdi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	DRM_IOCTL_BASE   uintptr = 'd'
	DRM_COMMAND_BASE uintptr = 0x40
)

func DRM_IO(nr uintptr) uintptr {
	return IO(DRM_IOCTL_BASE, nr)
}

func DRM_IOW(nr, size uintptr) uintptr {
	return IOW(DRM_IOCTL_BASE, nr, size)
}

func DRM_IOWR(nr, size uintptr) uintptr {
	return IOWR(DRM_IOCTL_BASE, nr, size)
}

var (
	DRM_IOCTL_VERSION = DRM_IOWR(0x00, unsafe.Sizeof(drmVersion{}))

	DRM_IOCTL_AUTH_MAGIC  = DRM_IOW(0x11, unsafe.Sizeof(drmAuth{}))
	DRM_IOCTL_DROP_MASTER = DRM_IO(0x1f)

	DRM_IOCTL_EVDI_CONNECT = DRM_IOWR(DRM_COMMAND_BASE+DRM_EVDI_CONNECT, unsafe.Sizeof(drmEvdiConnect{}))
	// DRM_IOCTL_EVDI_REQUEST_UPDATE
	// DRM_IOCTL_EVDI_GRABPIX
	// DRM_IOCTL_EVDI_DDCCI_RESPONSE
	DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS = DRM_IOWR(DRM_COMMAND_BASE+DRM_EVDI_ENABLE_CURSOR_EVENTS, unsafe.Sizeof(drmEvdiEnableCursorEvents{}))
)

const (
	/* Input ioctls from evdi lib to driver */
	DRM_EVDI_CONNECT              = 0x00
	DRM_EVDI_REQUEST_UPDATE       = 0x01
	DRM_EVDI_GRABPIX              = 0x02
	DRM_EVDI_DDCCI_RESPONSE       = 0x03
	DRM_EVDI_ENABLE_CURSOR_EVENTS = 0x04
)

type drmMagic uint32

type drmAuth struct {
	magic drmMagic
}

func ioctl(fd int, req uint, arg uintptr) (n int, err error) {
	r1, _, e1 := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), uintptr(arg))
	n = int(r1)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

func drmIoctl(fd uintptr, req uintptr, arg uintptr) syscall.Errno {
	var ret uintptr
	var errno syscall.Errno
	for {
		ret, _, errno = unix.Syscall(unix.SYS_IOCTL, fd, req, arg)
		// (ret == -1 && (errno == EINTR || erno == EAGAIN))
		do := (^ret == 0 && (errno == syscall.EINTR || errno == syscall.EAGAIN))
		if !do {
			break
		}
	}
	if errno != 0 {
		evdiLogDebug("drmIoctl: err=%v", errno)
	}
	return errno
}

func drmAuthMagic(fd uintptr, magic drmMagic) syscall.Errno {
	var auth drmAuth

	auth.magic = magic

	if ret := drmIoctl(fd, DRM_IOCTL_AUTH_MAGIC, uintptr(unsafe.Pointer(&auth))); ret != 0 {
		return ret
	}
	return 0
}

func drmIsMaster(fd uintptr) bool {
	res := drmAuthMagic(fd, 0)
	// e := ^(uintptr(syscall.EACCES))
	// return res != e
	return res != syscall.EACCES
}

func doIoctl(fd uintptr, request uintptr, data uintptr, msg string) syscall.Errno {
	err := drmIoctl(fd, request, data)
	if err != 0 {
		evdiLogDebug("ioctl %s error: %s", msg, unix.ErrnoName(err))
	}
	return err
}

type drmVersion struct {
	versionMajor      int32
	versionMinor      int32
	versionPatchLevel int32
	nameLen           uint32
	name              unsafe.Pointer
	dateLen           uint32
	date              unsafe.Pointer
	descLen           uint32
	desc              unsafe.Pointer
}

type drmEvdiConnect struct {
	connected           int32
	devIndex            int32
	edid                *byte
	edidLength          uint32
	pixelAreaLimit      uint32
	pixelPerSecondLimit uint32
}

type drmEvdiEnableCursorEvents struct {
	base   drmEvent
	enable uint8
}

type drmEvent struct {
	typ    uint32
	length uint32
}

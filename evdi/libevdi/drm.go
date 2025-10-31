//go:build linux

package libevdi

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/maitredede/puregolibs/drm"
	ioc "github.com/maitredede/puregolibs/drm/ioctl"
	"golang.org/x/sys/unix"
)

const (
	/* Input ioctls from evdi lib to driver */
	DRM_EVDI_CONNECT              = 0x00
	DRM_EVDI_REQUEST_UPDATE       = 0x01
	DRM_EVDI_GRABPIX              = 0x02
	DRM_EVDI_DDCCI_RESPONSE       = 0x03
	DRM_EVDI_ENABLE_CURSOR_EVENTS = 0x04

	DDCCI_BUFFER_SIZE = 64
)

var (
	DRM_IOCTL_EVDI_CONNECT              = ioc.NewCode(ioc.Read|ioc.Write, uint16(unsafe.Sizeof(drmEvdiConnect{})), drm.IOCTLBase, drm.CommandBase+DRM_EVDI_CONNECT)
	DRM_IOCTL_EVDI_REQUEST_UPDATE       = ioc.NewCode(ioc.Read|ioc.Write, uint16(unsafe.Sizeof(drmEvdiRequestUpdate{})), drm.IOCTLBase, drm.CommandBase+DRM_EVDI_REQUEST_UPDATE)
	DRM_IOCTL_EVDI_GRABPIX              = ioc.NewCode(ioc.Read|ioc.Write, uint16(unsafe.Sizeof(drmEvdiGrabPix{})), drm.IOCTLBase, drm.CommandBase+DRM_EVDI_GRABPIX)
	DRM_IOCTL_EVDI_DDCCI_RESPONSE       = ioc.NewCode(ioc.Read|ioc.Write, uint16(unsafe.Sizeof(drmEvdiDdcciResponse{})), drm.IOCTLBase, drm.CommandBase+DRM_EVDI_DDCCI_RESPONSE)
	DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS = ioc.NewCode(ioc.Read|ioc.Write, uint16(unsafe.Sizeof(drmEvdiEnableCursorEvents{})), drm.IOCTLBase, drm.CommandBase+DRM_EVDI_ENABLE_CURSOR_EVENTS)
)

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

func doIoctl(f *os.File, request uint32, data uintptr, msg string) syscall.Errno {
	err := drmIoctl(f.Fd(), uintptr(request), data)
	if err != 0 {
		evdiLogDebug("ioctl %s error: %s", msg, unix.ErrnoName(err))
	}
	return err
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
	typ    DrmEvdiEventType
	length uint32
}

type drmEvdiRequestUpdate struct {
	reserved int32
}

type grabPixMode int32

const (
	EVDI_GRABPIX_MODE_RECTS grabPixMode = 0
	EVDI_GRABPIX_MODE_DIRTY grabPixMode = 1
)

type drmEvdiGrabPix struct {
	mode          grabPixMode
	bufWidth      int32
	bufHeight     int32
	bufByteStride int32
	buffer        *byte
	numRects      int32
	rects         *drmClipRect
}

type drmClipRect struct {
	x1, y1, x2, y2 uint16
}

type drmEvdiDdcciResponse struct {
	buffer       *byte
	bufferLength uint32
	result       uint8
}

type drmEvdiEventDdcciData struct {
	base         drmEvent
	buffer       [DDCCI_BUFFER_SIZE]byte
	bufferLength uint32
	flags        uint16
	address      uint16
}

type drmEvdiEventDpms struct {
	base drmEvent
	mode int32
}

type drmEvdiModeChanged struct {
	base         drmEvent
	hdisplay     int32
	vdisplay     int32
	vrefresh     int32
	bitsPerPixel int32
	pixelFormat  int32
}

type drmEvdiEventCrtcState struct {
	base  drmEvent
	state int32
}

type drmEvdiEventCursorSet struct {
	base         drmEvent
	hotX         int32
	hotY         int32
	width        uint32
	height       uint32
	enabled      uint8
	bufferHandle uint32
	bufferLength uint32
	pixelFormat  uint32
	stride       uint32
}

type drmEvdiEventCursorMove struct {
	base drmEvent
	x    int32
	y    int32
}

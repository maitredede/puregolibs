//go:build linux

package libevdi

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"

	drmmode "github.com/maitredede/puregolibs/drm/mode"
	"github.com/maitredede/puregolibs/tools/linkedlist"
	"golang.org/x/sys/unix"
)

type Handle struct {
	fd          *os.File
	deviceIndex int

	buffersMap           map[int32]bufferData
	frameBuffersListHead linkedlist.LinkedList[*evdiBuffer]
	bufferToUpdate       int32
}

const (
	EvdiInvalidDeviceIndex = -1
	EvdiUsageLength        = 64

	bufferCount = 2
)

func OpenAttachedTo(sysfsParentDevice string) (*Handle, error) {
	deviceIndex := EvdiInvalidDeviceIndex
	if len(sysfsParentDevice) == 0 {
		deviceIndex = getGenericDevice()
	} else {
		if strings.HasPrefix(sysfsParentDevice, "usb:") {
			deviceIndex = getDeviceAttachedToUsb(sysfsParentDevice)
		} else {
			return nil, errors.New("unrecognized parent identifier")
		}
	}

	if deviceIndex >= 0 && deviceIndex < EvdiUsageLength {
		return Open(deviceIndex)
	}
	return nil, errors.New("open failed")
}

func Open(device int) (*Handle, error) {
	fd, err := openDevice(device)
	if err != nil {
		return nil, fmt.Errorf("device open failed: %w", err)
	}

	if !isEvdi(fd) {
		fd.Close()
		return nil, fmt.Errorf("device is not evdi")
	}
	if !isEvdiCompatible(fd) {
		fd.Close()
		return nil, fmt.Errorf("device is not evdi-compatible")
	}

	h := &Handle{
		fd:          fd,
		deviceIndex: device,

		buffersMap: make(map[int32]bufferData),
	}
	cardUsage[device] = h
	evdiLogInfo("using /dev/dri/card%d", device)
	return h, nil
}

func (h *Handle) GetEventReady() uintptr {
	return h.fd.Fd()
}

func (h *Handle) Close() error {
	var errs []error
	errs = append(errs, h.fd.Close())

	for i, elem := range cardUsage {
		if elem == h {
			cardUsage[i] = nil
			evdiLogInfo("Marking /dev/dri/card%d as unused", h.deviceIndex)
		}
	}
	return errors.Join(errs...)
}

func (h *Handle) Connect(edid []byte, skuAreaLimit uint32) {
	h.Connect2(edid, skuAreaLimit, skuAreaLimit*60)
}

func (h *Handle) Connect2(edid []byte, pixelAreaLimit uint32, pixelPerSecondLimit uint32) {
	cmd := drmEvdiConnect{
		connected:           1,
		devIndex:            int32(h.deviceIndex),
		edid:                &edid[0],
		edidLength:          uint32(len(edid)),
		pixelAreaLimit:      pixelAreaLimit,
		pixelPerSecondLimit: pixelPerSecondLimit,
	}
	doIoctl(h.fd, DRM_IOCTL_EVDI_CONNECT, uintptr(unsafe.Pointer(&cmd)), "connect")
}

func (h *Handle) Disconnect() {
	cmd := drmEvdiConnect{}

	doIoctl(h.fd, DRM_IOCTL_EVDI_CONNECT, uintptr(unsafe.Pointer(&cmd)), "disconnect")
}

func (h *Handle) EnableCursorEvents(enabled bool) {
	cmd := drmEvdiEnableCursorEvents{
		enable: 0,
	}
	msg := "disabling"
	if enabled {
		cmd.enable = 1
		msg = "enabling"
	}

	evdiLogInfo("%s events on /dev/dri/card%d", msg, h.deviceIndex)
	err := doIoctl(h.fd, DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS, uintptr(unsafe.Pointer(&cmd)), "enable cursor events")
	if err != 0 {
		evdiLogDebug("cursor events ret: %v", err)
	}
}

func (h *Handle) PollEvents(timeoutMS int, handlers EventHandlers) {
	fd := h.fd.Fd()
	events := []unix.PollFd{
		{Fd: int32(fd), Events: unix.POLLIN},
	}
	for {
		n, err := unix.Poll(events, timeoutMS)
		if err != nil {
			if err == unix.EINTR {
				continue
			}
			evdiLogDebug("epoll wait error: %v\n", err)
			continue
		}
		if n != 0 {
			h.HandleEvents(handlers)
		}
		break
	}
}

func (h *Handle) HandleEvents(handlers EventHandlers) {
	buffer := make([]byte, 1024)
	bytesRead, err := h.fd.Read(buffer)
	if err != nil {
		evdiLogDebug("TODO handleEvents read error: %v", err)
	}
	var i int
	for {
		if i >= bytesRead {
			break
		}
		e := (*drmEvent)(unsafe.Pointer(&buffer[i]))
		evdiLogDebug("events: i=%d/%d e.typ=%s e.length=%d", i, bytesRead, e.typ.String(), e.length)
		h.dispatchEvent(handlers, e)

		i += int(e.length)
	}
}

func (h *Handle) dispatchEvent(evtctx EventHandlers, e *drmEvent) {
	switch e.typ {
	case DRM_EVDI_EVENT_UPDATE_READY:
		if evtctx.UpdateReady != nil {
			evtctx.UpdateReady(h.bufferToUpdate, evtctx.UserData)
		}
	case DRM_EVDI_EVENT_DPMS:
		if evtctx.Dpms != nil {
			event := *(*drmEvdiEventDpms)(unsafe.Pointer(e))
			evtctx.Dpms(DpmsMode(event.mode), evtctx.UserData)
		}
	case DRM_EVDI_EVENT_MODE_CHANGED:
		if evtctx.ModeChanged != nil {
			event := *(*drmEvdiModeChanged)(unsafe.Pointer(e))
			evtctx.ModeChanged(toEvdiMode(event), evtctx.UserData)
		}
	case DRM_EVDI_EVENT_CRTC_STATE:
		if evtctx.CrtcState != nil {
			event := *(*drmEvdiEventCrtcState)(unsafe.Pointer(e))
			evtctx.CrtcState(event.state, evtctx.UserData)
		}
	case DRM_EVDI_EVENT_CURSOR_SET:
		if evtctx.CursorSet != nil {
			event := *(*drmEvdiEventCursorSet)(unsafe.Pointer(e))
			cursorSet := h.toEvdiCursorSet(event)
			if cursorSet.Enabled && len(cursorSet.BufferData) == 0 {
				evdiLogInfo("Error: Cursor buffer is null!")
				evdiLogInfo("Disabling cursor events")
				h.EnableCursorEvents(false)
				cursorSet.Enabled = false
				cursorSet.BufferData = nil
			}
			evtctx.CursorSet(cursorSet, evtctx.UserData)
		}
	case DRM_EVDI_EVENT_CURSOR_MOVE:
		if evtctx.CursorMove != nil {
			event := *(*drmEvdiEventCursorMove)(unsafe.Pointer(e))
			evtctx.CursorMove(toEvdiCursorMove(event), evtctx.UserData)
		}
	case DRM_EVDI_EVENT_DDCCI_DATA:
		if evtctx.DdcciData != nil {
			event := *(*drmEvdiEventDdcciData)(unsafe.Pointer(e))
			evtctx.DdcciData(toEvdiDdcciData(event), evtctx.UserData)
		}
	default:
		evdiLogInfo("warning: unhandled event 0x%08x", e.typ)
	}
}

func (h *Handle) UnregisterBuffer(bufferID int32) {
	// entry := h.findBuffer(bufferID)
	// if entry == nil {
	// 	return
	// }
	h.removeFrameBuffer(bufferID)
}

// func (h *Handle) findBuffer(id int32) *evdiBuffer {
// 	node := h.frameBuffersListHead.FindItem(func(a *evdiBuffer) bool {
// 		return a.id == id
// 	})
// 	if node != nil {
// 		return node.Item
// 	}
// 	return nil
// }

func (h *Handle) removeFrameBuffer(bufferID int32) {
	entry := h.frameBuffersListHead.FindItem(func(a *evdiBuffer) bool {
		return a.id == bufferID
	})
	if entry == nil {
		return
	}
	// TODO check buffer alloc/free
	// libc.Free(unsafe.Pointer(entry.Item.buffer))
}

func (h *Handle) RegisterBuffer(buffer *evdiBuffer) {
	entry := h.frameBuffersListHead.FindItem(func(a *evdiBuffer) bool { return a.id == buffer.id })
	if entry != nil {
		panic("buffer already exists")
	}
	h.addFrameBuffer(buffer)
}

func (h *Handle) addFrameBuffer(buffer *evdiBuffer) {
	h.frameBuffersListHead.Append(buffer)
}

func (h *Handle) RequestUpdate(bufferID int32) bool {
	evdiLogDebug("called requestUpdate id=%d", bufferID)
	h.bufferToUpdate = bufferID

	cmd := drmEvdiRequestUpdate{}
	requestResult := doIoctl(h.fd, DRM_IOCTL_EVDI_REQUEST_UPDATE, uintptr(unsafe.Pointer(&cmd)), "request_update")
	grabImmediately := requestResult == 1
	return grabImmediately
}

func (h *Handle) getDumbOffset(handle uint32, offset *uint64) error {
	// mapDumb := drmModeMapDumb{}
	// mapDumb.handle = handle
	// ret := doIoctl(h.fd, DRM_IOCTL_MODE_MAP_DUMB, uintptr(unsafe.Pointer(&mapDumb)), "DRM_MODE_MAP_DUMB")
	// *offset = mapDumb.offset
	// return ret

	res, err := drmmode.MapDumb(h.fd, handle)
	if err != nil {
		evdiLogDebug("handle: getDumbOffset err=%v", err)
	}
	*offset = res
	return err
}

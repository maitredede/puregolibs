package libevdi

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

type Handle struct {
	fd          *os.File
	deviceIndex int
}

const (
	EvdiInvalidDeviceIndex = -1
	EvdiUsageLength        = 64
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
	doIoctl(h.fd, DRM_IOCTL_EVDI_ENABLE_CURSOR_EVENTS, uintptr(unsafe.Pointer(&cmd)), "enable cursor events")
}

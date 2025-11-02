//go:build linux

package libevdi

import "fmt"

type DrmEvdiEventType uint32

const (
	DRM_EVDI_EVENT_UPDATE_READY DrmEvdiEventType = 0x80000000
	DRM_EVDI_EVENT_DPMS         DrmEvdiEventType = 0x80000001
	DRM_EVDI_EVENT_MODE_CHANGED DrmEvdiEventType = 0x80000002
	DRM_EVDI_EVENT_CRTC_STATE   DrmEvdiEventType = 0x80000003
	DRM_EVDI_EVENT_CURSOR_SET   DrmEvdiEventType = 0x80000004
	DRM_EVDI_EVENT_CURSOR_MOVE  DrmEvdiEventType = 0x80000005
	DRM_EVDI_EVENT_DDCCI_DATA   DrmEvdiEventType = 0x80000006
)

func (t DrmEvdiEventType) String() string {
	switch t {
	case DRM_EVDI_EVENT_UPDATE_READY:
		return "DRM_EVDI_EVENT_UPDATE_READY"
	case DRM_EVDI_EVENT_DPMS:
		return "DRM_EVDI_EVENT_DPMS"
	case DRM_EVDI_EVENT_MODE_CHANGED:
		return "DRM_EVDI_EVENT_MODE_CHANGED"
	case DRM_EVDI_EVENT_CRTC_STATE:
		return "DRM_EVDI_EVENT_CRTC_STATE"
	case DRM_EVDI_EVENT_CURSOR_SET:
		return "DRM_EVDI_EVENT_CURSOR_SET"
	case DRM_EVDI_EVENT_CURSOR_MOVE:
		return "DRM_EVDI_EVENT_CURSOR_MOVE"
	case DRM_EVDI_EVENT_DDCCI_DATA:
		return "DRM_EVDI_EVENT_DDCCI_DATA"
	}
	return fmt.Sprintf("DRM_EVDI_EVENT_???_%08x", uint32(t))
}

type EventHandlers struct {
	UserData any

	UpdateReady func(bufferToUpdate int32, userData any)
	Dpms        func(dpms DpmsMode, userData any)
	ModeChanged func(mode Mode, userData any)
	CrtcState   func(state int32, userData any)
	CursorSet   func(cursor CursorSet, userData any)
	CursorMove  func(cursor CursorMove, userData any)
	DdcciData   func(data DdcciData, userData any)
}

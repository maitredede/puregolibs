package libevdi

import "fmt"

type DpmsMode int32

const (
	DPMSModeOn      DpmsMode = 0
	DPMSModeStandby DpmsMode = 1
	DPMSModeSuspend DpmsMode = 2
	DPMSModeOff     DpmsMode = 3
)

func (m DpmsMode) String() string {
	switch m {
	case DPMSModeOn:
		return "on"
	case DPMSModeStandby:
		return "standby"
	case DPMSModeSuspend:
		return "suspend"
	case DPMSModeOff:
		return "off"
	default:
		return fmt.Sprintf("%x", int32(m))
	}
}

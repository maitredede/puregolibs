package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type PowerStatus int32

const (
	PowerStatusOn                      PowerStatus = 0x00
	PowerStatusStandby                 PowerStatus = 0x01
	PowerStatusInTransitionStandbyToOn PowerStatus = 0x02
	PowerStatusInTransitionOnToStandby PowerStatus = 0x03
	PowerStatusUnknown                 PowerStatus = 0x99
)

func (s PowerStatus) String() string {
	libInit()

	buffSize := int32(1024)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecPowerStatusToString(s, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}

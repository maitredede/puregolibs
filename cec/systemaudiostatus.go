package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type SystemAudioStatus int32

const (
	SystemAudioStatusOn  SystemAudioStatus = 0
	SystemAudioStatusOff SystemAudioStatus = 1
)

func (s SystemAudioStatus) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecSystemAudioStatusToString(s, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}

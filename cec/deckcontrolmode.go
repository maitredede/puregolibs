package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type DeckControlMode int32

const (
	DeckControlModeSkipForwardWind   DeckControlMode = 1
	DeckControlModeSkipReverseRewind DeckControlMode = 2
	DeckControlModeStop              DeckControlMode = 3
	DeckControlModeEject             DeckControlMode = 4
)

func (m DeckControlMode) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecDeckControlModeToString(m, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}

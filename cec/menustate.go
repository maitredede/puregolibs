package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type MenuState int32

const (
	MenuStateActivated   MenuState = 0
	MenuStateDeactivated MenuState = 1
)

func (s MenuState) String() string {
	libInit()

	buffSize := int32(1024)
	buff := make([]byte, buffSize)
	buffPtr := unsafe.Pointer(&buff[0])
	libCecMenuStateToString(s, buffPtr, buffSize)
	return strings.GoStringN((*byte)(buffPtr), int(buffSize))
}

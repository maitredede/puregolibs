package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
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
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecMenuStateToString(s, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}

package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type CecVersion int32

const (
	CecVersionUnknown CecVersion = 0x00
	CecVersion_1_2    CecVersion = 0x01
	CecVersion_1_2A   CecVersion = 0x02
	CecVersion_1_3    CecVersion = 0x03
	CecVersion_1_3A   CecVersion = 0x04
	CecVersion_1_4    CecVersion = 0x05
	CecVersion_2_0    CecVersion = 0x06
)

func (v CecVersion) String() string {
	libInit()

	buffSize := int32(128)
	buff := make([]byte, buffSize)
	buffPtr := unsafe.Pointer(&buff[0])
	libCecCecVersionToString(v, buffPtr, buffSize)
	return strings.GoStringN((*byte)(buffPtr), int(buffSize))
}

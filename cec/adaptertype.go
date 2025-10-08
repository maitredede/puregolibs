package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

type AdapterType int32

const (
	AdapterTypeUnknown         AdapterType = 0
	AdapterTypeP8External      AdapterType = 0x1
	AdapterTypeP8DaughterBoard AdapterType = 0x2
	AdapterTypeRPI             AdapterType = 0x100
	AdapterTypeTDA995x         AdapterType = 0x200
	AdapterTypeExynos          AdapterType = 0x300
	AdapterTypeLinux           AdapterType = 0x400
	AdapterTypeAOCEC           AdapterType = 0x500
	AdapterTypeIMX             AdapterType = 0x600
)

func (t AdapterType) String() string {
	libInit()

	buffSize := int32(1024)
	buff := make([]byte, buffSize)
	buffPtr := uintptr(unsafe.Pointer(&buff[0]))
	libCecAdapterTypeToString(t, buffPtr, buffSize)
	return strings.GoStringN(buffPtr, int(buffSize))
}

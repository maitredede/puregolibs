package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
)

type LogicalAddress int32

const (
	DeviceUnknown LogicalAddress = iota - 1
	DeviceTV
	DeviceRecordingDevice1
	DeviceRecordingDevice2
	DeviceTuner1
	DevicePlaybackDevice1
	DeviceAudioSystem
	DeviceTuner2
	DeviceTuner3
	DevicePlaybackDevice2
	DeviceRecordingDevice3
	DeviceTuner4
	DevicePlaybackDevice3
	DeviceReserved1
	DeviceReserved2
	DeviceFreeUse
	DeviceUnregistered
	DeviceBroadcast LogicalAddress = 15
)

const logicalAddressesCount = 16

func (a LogicalAddress) String() string {
	libInit()

	buffSize := int32(1024)
	buff := make([]byte, buffSize)
	buffPtr := unsafe.Pointer(&buff[0])
	libCecLogicalAddressToString(a, buffPtr, buffSize)
	return strings.GoStringN((*byte)(buffPtr), int(buffSize))
}

type LogicalAddresses struct {
	primary   LogicalAddress
	addresses [logicalAddressesCount]int32
}

func NewLogicalAddresses() *LogicalAddresses {
	l := &LogicalAddresses{}
	l.Clear()
	return l
}

func (l *LogicalAddresses) Clear() {
	l.primary = DeviceUnregistered
	for i := 0; i < logicalAddressesCount; i++ {
		l.addresses[i] = 0
	}
}

func (l *LogicalAddresses) IsEmpty() bool {
	return l.primary == DeviceUnregistered
}

func (l *LogicalAddresses) Set(address LogicalAddress) {
	if l.primary == DeviceUnregistered {
		l.primary = address
	}
	l.addresses[int(address)] = 1
}

func (l *LogicalAddresses) Unset(address LogicalAddress) {
	if l.primary == address {
		l.primary = DeviceUnregistered
	}
	l.addresses[int(address)] = 0
}

func (l *LogicalAddresses) IsSet(address LogicalAddress) bool {
	return l.addresses[int(address)] == 1
}

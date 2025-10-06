package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)


type Configuration struct {
	ClientVersion   VersionValue
	DeviceName      string
	DeviceTypes     DeviceTypeList
	AutodectAddress bool
	PhysicalAddress int
	BaseDevice      LogicalAddress
	HDMIPort        byte
	TVVendor        uint32
	WakeDevices     LogicalAddresses
	PowerDevices    LogicalAddresses

	ServerVersion VersionValue

	GetSettingsFromROM bool
	ActivateSource     bool
	PowerOffOnStandby  bool
	//TODO callbacks
	LogicalAddress  LogicalAddress
	FirmwareVersion uint16
}

type DeviceNameString [OSDNameSize]byte

func (s DeviceNameString) String() string {
	return strings.GoString(uintptr(unsafe.Pointer(&s[0])))
}

type NativeConfiguration struct {
	ClientVersion   VersionValue
	DeviceName      DeviceNameString
	DeviceTypes     DeviceTypeList
	AutodectAddress byte
	PhysicalAddress uint16
	BaseDevice      LogicalAddress
	HDMIPort        byte
	TVVendor        uint32
	WakeDevices     LogicalAddresses
	PowerDevices    LogicalAddresses

	ServerVersion VersionValue

	GetSettingsFromROM byte
	ActivateSource     byte
	PowerOffOnStandby  byte

	CallbackParam uintptr
	Callbacks     uintptr

	LogicalAddress    LogicalAddress
	FirmwareVersion   uint16
	DeviceLanguage    [3]byte
	FirmwareBuildDate uint32
	MonitorOnly       byte
	CecVersion        VersionValue
	AdapterType       AdapterType

	//TODO
	padding [88]byte
}

func toBool(b bool) byte {
	if b {
		return 1
	}
	return 0
}

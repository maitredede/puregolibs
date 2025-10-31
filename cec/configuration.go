package cec

import (
	"unsafe"

	"github.com/maitredede/puregolibs/tools/strings"
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
	p := &s[0]
	return strings.GoStringN(p, OSDNameSize)
}

func CDeviceNameString(name string) DeviceNameString {
	bin := []byte(name)
	var value DeviceNameString

	max := min(OSDNameSize, len(bin))
	for i := 0; i < max; i++ {
		value[i] = bin[i]
	}
	return value
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

	// player specific settings
	GetSettingsFromROM byte
	ActivateSource     byte
	PowerOffOnStandby  byte

	CallbackParam unsafe.Pointer
	//Callbacks     unsafe.Pointer
	Callbacks *nativeICECCallbacks

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

func (n NativeConfiguration) Go() Configuration {
	cfg := Configuration{
		ClientVersion:      n.ClientVersion,
		DeviceName:         n.DeviceName.String(),
		DeviceTypes:        n.DeviceTypes,
		AutodectAddress:    n.AutodectAddress != 0,
		PhysicalAddress:    int(n.PhysicalAddress),
		BaseDevice:         n.BaseDevice,
		HDMIPort:           n.HDMIPort,
		TVVendor:           n.TVVendor,
		WakeDevices:        n.WakeDevices,
		PowerDevices:       n.PowerDevices,
		ServerVersion:      n.ServerVersion,
		GetSettingsFromROM: n.GetSettingsFromROM != 0,
		ActivateSource:     n.ActivateSource != 0,
		PowerOffOnStandby:  n.PowerOffOnStandby != 0,
		LogicalAddress:     n.LogicalAddress,
		FirmwareVersion:    n.FirmwareVersion,
		//TODO
	}
	return cfg
}

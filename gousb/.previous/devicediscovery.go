package gousb

// DeviceDiscovery controls USB device discovery.
type DeviceDiscovery int

const (
	// EnableDeviceDiscovery means the connected USB devices will be enumerated
	// on Context initialization. This enables the use of OpenDevices and
	// OpenWithVIDPID. This is the default.
	EnableDeviceDiscovery DeviceDiscovery = iota
	// DisableDeviceDiscovery means the USB devices are not enumerated and
	// OpenDevices will not return any devices.
	// Without device discovery, OpenDeviceWithFileDescriptor can be used
	// to open devices.
	DisableDeviceDiscovery
)

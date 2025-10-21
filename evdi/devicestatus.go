package evdi

type DeviceStatus int32

const (
	DeviceStatusAvailable DeviceStatus = iota
	DeviceStatusUnrecognized
	DeviceStatusNotPresent
)

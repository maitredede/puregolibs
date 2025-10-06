package cec

type DeviceTypeList [5]DeviceType

type DeviceType int32

const (
	DeviceTypeTV DeviceType = iota
	DeviceTypeRecordingDevice
	DeviceTypeReserved
	DeviceTypeTuner
	DeviceTypePlaybackDevice
	DeviceTypeAudioSystem
)

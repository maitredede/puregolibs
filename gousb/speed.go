package gousb

// Speed identifies the speed of the device.
type Speed int

// Device speeds as defined in the USB spec.
const (
	SpeedUnknown     Speed = 0 // The OS doesn't report or know the device speed
	SpeedLow         Speed = 1 // The device is operating at low speed (1.5MBit/s)
	SpeedFull        Speed = 2 // The device is operating at full speed (12MBit/s)
	SpeedHigh        Speed = 3 // The device is operating at high speed (480MBit/s)
	SpeedSuper       Speed = 4 // The device is operating at super speed (5000MBit/s)
	SpeedSuperPlus   Speed = 5 // The device is operating at super speed plus (10000MBit/s)
	SpeedSuperPlusX2 Speed = 6 // The device is operating at super speed plus x2 (20000MBit/s)
)

var deviceSpeedDescription = map[Speed]string{
	SpeedUnknown:     "unknown",
	SpeedLow:         "low",
	SpeedFull:        "full",
	SpeedHigh:        "high",
	SpeedSuper:       "super",
	SpeedSuperPlus:   "superPlus",
	SpeedSuperPlusX2: "superPlusX2",
}

// String returns a human-readable name of the device speed.
func (s Speed) String() string {
	return deviceSpeedDescription[s]
}

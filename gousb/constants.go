package gousb

const (
	selfPoweredMask  = 0x40
	remoteWakeupMask = 0x20
)

// Milliamperes is a unit of electric current consumption.
type Milliamperes uint

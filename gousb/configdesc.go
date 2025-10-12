package gousb

import "fmt"

// ConfigDesc contains the information about a USB device configuration,
// extracted from the device descriptor.
type ConfigDesc struct {
	// Number is the configuration number.
	Number int
	// SelfPowered is true if the device is powered externally, i.e. not
	// drawing power from the USB bus.
	SelfPowered bool
	// RemoteWakeup is true if the device supports remote wakeup, i.e.
	// an external signal that will wake up a suspended USB device. An example
	// might be a keyboard that can wake up through a keypress after
	// the host put it in suspend mode. Note that gousb does not support
	// device power management, RemoteWakeup only refers to the reported device
	// capability.
	RemoteWakeup bool
	// MaxPower is the maximum current the device draws from the USB bus
	// in this configuration.
	MaxPower Milliamperes
	// Interfaces has a list of USB interfaces available in this configuration.
	Interfaces []InterfaceDesc

	iConfiguration int // index of a string descriptor describing this configuration
}

// String returns the human-readable description of the configuration descriptor.
func (c ConfigDesc) String() string {
	return fmt.Sprintf("Configuration %d", c.Number)
}

func (c ConfigDesc) intfDesc(num int) (*InterfaceDesc, error) {
	// In an ideal world, interfaces in the descriptor would be numbered
	// contiguously starting from 0, as required by the specification. In the
	// real world however the specification is sometimes ignored:
	// https://github.com/google/gousb/issues/65
	ifs := make([]int, len(c.Interfaces))
	for i, desc := range c.Interfaces {
		if desc.Number == num {
			return &desc, nil
		}
		ifs[i] = desc.Number
	}
	return nil, fmt.Errorf("interface %d not found, available interface numbers: %v", num, ifs)
}

package gousb

import (
	"fmt"
	"sort"
)

// DeviceDesc is a representation of a USB device descriptor.
type DeviceDesc struct {
	// The bus on which the device was detected
	Bus int
	// The address of the device on the bus
	Address int
	// The negotiated operating speed for the device
	Speed Speed
	// The pyhsical port on the parent hub on which the device is connected.
	// Ports are numbered from 1, excepting root hub devices which are always 0.
	Port int
	// Physical path of connected parent ports, starting at the root hub device.
	// A path length of 0 represents a root hub device,
	// a path length of 1 represents a device directly connected to a root hub,
	// a path length of 2 or more are connected to intermediate hub devices.
	// e.g. [1,2,3] represents a device connected to port 3 of a hub connected
	// to port 2 of a hub connected to port 1 of a root hub.
	Path []int

	// Version information
	Spec   BCD // USB Specification Release Number
	Device BCD // The device version

	// Product information
	Vendor  ID // The Vendor identifer
	Product ID // The Product identifier

	// Protocol information
	Class                Class    // The class of this device
	SubClass             Class    // The sub-class (within the class) of this device
	Protocol             Protocol // The protocol (within the sub-class) of this device
	MaxControlPacketSize int      // Maximum size of the control transfer

	// Configuration information
	Configs map[int]ConfigDesc

	iManufacturer int // The Manufacturer descriptor index
	iProduct      int // The Product descriptor index
	iSerialNumber int // The SerialNumber descriptor index
}

// String returns a human-readable version of the device descriptor.
func (d *DeviceDesc) String() string {
	return fmt.Sprintf("%d.%d: %s:%s (available configs: %v)", d.Bus, d.Address, d.Vendor, d.Product, d.sortedConfigIds())
}

func (d *DeviceDesc) sortedConfigIds() []int {
	var cfgs []int
	for c := range d.Configs {
		cfgs = append(cfgs, c)
	}
	sort.Ints(cfgs)
	return cfgs
}

func (d *DeviceDesc) cfgDesc(cfgNum int) (*ConfigDesc, error) {
	desc, ok := d.Configs[cfgNum]
	if !ok {
		return nil, fmt.Errorf("configuration id %d not found in the descriptor of the device. Available config ids: %v", cfgNum, d.sortedConfigIds())
	}
	return &desc, nil
}

package gousb

import (
	"fmt"
	"sort"
)

type libusbDeviceDescriptor struct {
	/** Size of this descriptor (in bytes) */
	bLength uint8

	/** Descriptor type. Will have value
	 * \ref libusb_descriptor_type::LIBUSB_DT_DEVICE LIBUSB_DT_DEVICE in this
	 * context. */
	bDescriptorType uint8

	/** USB specification release number in binary-coded decimal. A value of
	 * 0x0200 indicates USB 2.0, 0x0110 indicates USB 1.1, etc. */
	bcdUSB uint16

	/** USB-IF class code for the device. See \ref libusb_class_code. */
	bDeviceClass uint8

	/** USB-IF subclass code for the device, qualified by the bDeviceClass
	 * value */
	bDeviceSubClass uint8

	/** USB-IF protocol code for the device, qualified by the bDeviceClass and
	 * bDeviceSubClass values */
	bDeviceProtocol uint8

	/** Maximum packet size for endpoint 0 */
	bMaxPacketSize0 uint8

	/** USB-IF vendor ID */
	idVendor uint16

	/** USB-IF product ID */
	idProduct uint16

	/** Device release number in binary-coded decimal */
	bcdDevice uint16

	/** Index of string descriptor describing manufacturer */
	iManufacturer uint8

	/** Index of string descriptor describing product */
	iProduct uint8

	/** Index of string descriptor containing device serial number */
	iSerialNumber uint8

	/** Number of possible configurations */
	bNumConfigurations uint8
}

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

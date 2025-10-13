package gousb

import (
	"fmt"
	"unsafe"
)

type libusbConfigDescriptor struct {
	/** Size of this descriptor (in bytes) */
	bLength uint8

	/** Descriptor type. Will have value
	 * \ref libusb_descriptor_type::LIBUSB_DT_CONFIG LIBUSB_DT_CONFIG
	 * in this context. */
	bDescriptorType uint8

	/** Total length of data returned for this configuration */
	wTotalLength uint16

	/** Number of interfaces supported by this configuration */
	bNumInterfaces uint8

	/** Identifier value for this configuration */
	bConfigurationValue uint8

	/** Index of string descriptor describing this configuration */
	iConfiguration uint8

	/** Configuration characteristics */
	bmAttributes uint8

	/** Maximum power consumption of the USB device from this bus in this
	 * configuration when the device is fully operation. Expressed in units
	 * of 2 mA when the device is operating in high-speed mode and in units
	 * of 8 mA when the device is operating in super-speed mode. */
	MaxPower uint8

	/** Array of interfaces supported by this configuration. The length of
	 * this array is determined by the bNumInterfaces field. */
	//const struct libusb_interface *interface;
	iface *libusbInterface

	/** Extra descriptors. If libusb encounters unknown configuration
	 * descriptors, it will store them here, should you wish to parse them. */
	//const unsigned char *extra;
	extra unsafe.Pointer

	/** Length of the extra descriptors, in bytes. Must be non-negative. */
	extra_length int32
}

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

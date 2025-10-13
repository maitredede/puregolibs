package gousb

import (
	"fmt"
	"unsafe"
)

type libusbInterfaceDescriptor struct {
	/** Size of this descriptor (in bytes) */
	bLength uint8

	/** Descriptor type. Will have value
	 * \ref libusb_descriptor_type::LIBUSB_DT_INTERFACE LIBUSB_DT_INTERFACE
	 * in this context. */
	bDescriptorType uint8

	/** Number of this interface */
	bInterfaceNumber uint8

	/** Value used to select this alternate setting for this interface */
	bAlternateSetting uint8

	/** Number of endpoints used by this interface (excluding the control
	 * endpoint). */
	bNumEndpoints uint8

	/** USB-IF class code for this interface. See \ref libusb_class_code. */
	bInterfaceClass uint8

	/** USB-IF subclass code for this interface, qualified by the
	 * bInterfaceClass value */
	bInterfaceSubClass uint8

	/** USB-IF protocol code for this interface, qualified by the
	 * bInterfaceClass and bInterfaceSubClass values */
	bInterfaceProtocol uint8

	/** Index of string descriptor describing this interface */
	iInterface uint8

	/** Array of endpoint descriptors. This length of this array is determined
	 * by the bNumEndpoints field. */
	//const struct libusb_endpoint_descriptor *endpoint;
	endpoint *libusbEndpointDescriptor

	/** Extra descriptors. If libusb encounters unknown interface descriptors,
	 * it will store them here, should you wish to parse them. */
	//const unsigned char *extra;
	extra unsafe.Pointer

	/** Length of the extra descriptors, in bytes. Must be non-negative. */
	extraLength int32
}

// InterfaceDesc contains information about a USB interface, extracted from
// the descriptor.
type InterfaceDesc struct {
	// Number is the number of this interface.
	Number int
	// AltSettings is a list of alternate settings supported by the interface.
	AltSettings []InterfaceSetting
}

func (i *InterfaceDesc) altSetting(alt int) (*InterfaceSetting, error) {
	alts := make([]int, len(i.AltSettings))
	for a, s := range i.AltSettings {
		if s.Alternate == alt {
			return &s, nil
		}
		alts[a] = s.Alternate
	}
	return nil, fmt.Errorf("alternate setting %d not found for %s, available alt settings: %v", alt, i, alts)
}

// String returns a human-readable description of the interface descriptor and
// its alternate settings.
func (i InterfaceDesc) String() string {
	return fmt.Sprintf("Interface %d (%d alternate settings)", i.Number, len(i.AltSettings))
}

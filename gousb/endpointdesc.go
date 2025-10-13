package gousb

import (
	"time"
	"unsafe"
)

type libusbEndpointDescriptor struct {
	/** Size of this descriptor (in bytes) */
	bLength uint8

	/** Descriptor type. Will have value
	 * \ref libusb_descriptor_type::LIBUSB_DT_ENDPOINT LIBUSB_DT_ENDPOINT in
	 * this context. */
	bDescriptorType uint8

	/** The address of the endpoint described by this descriptor. Bits 0:3 are
	 * the endpoint number. Bits 4:6 are reserved. Bit 7 indicates direction,
	 * see \ref libusb_endpoint_direction. */
	bEndpointAddress uint8

	/** Attributes which apply to the endpoint when it is configured using
	 * the bConfigurationValue. Bits 0:1 determine the transfer type and
	 * correspond to \ref libusb_endpoint_transfer_type. Bits 2:3 are only used
	 * for isochronous endpoints and correspond to \ref libusb_iso_sync_type.
	 * Bits 4:5 are also only used for isochronous endpoints and correspond to
	 * \ref libusb_iso_usage_type. Bits 6:7 are reserved. */
	bmAttributes uint8

	/** Maximum packet size this endpoint is capable of sending/receiving. */
	wMaxPacketSize uint16

	/** Interval for polling endpoint for data transfers. */
	bInterval uint8

	/** For audio devices only: the rate at which synchronization feedback
	 * is provided. */
	bRefresh uint8

	/** For audio devices only: the address if the synch endpoint */
	bSynchAddress uint8

	/** Extra descriptors. If libusb encounters unknown endpoint descriptors,
	 * it will store them here, should you wish to parse them. */
	//const unsigned char *extra;
	extra unsafe.Pointer

	/** Length of the extra descriptors, in bytes. Must be non-negative. */
	extraLength int32
}

// EndpointDesc contains the information about an interface endpoint, extracted
// from the descriptor.
type EndpointDesc struct {
	// Address is the unique identifier of the endpoint within the interface.
	Address EndpointAddress
	// Number represents the endpoint number. Note that the endpoint number is different from the
	// address field in the descriptor - address 0x82 means endpoint number 2,
	// with endpoint direction IN.
	// The device can have up to two endpoints with the same number but with
	// different directions.
	Number int
	// Direction defines whether the data is flowing IN or OUT from the host perspective.
	Direction EndpointDirection
	// MaxPacketSize is the maximum USB packet size for a single frame/microframe.
	MaxPacketSize int
	// TransferType defines the endpoint type - bulk, interrupt, isochronous.
	TransferType TransferType
	// PollInterval is the maximum time between transfers for interrupt and isochronous transfer,
	// or the NAK interval for a control transfer. See endpoint descriptor bInterval documentation
	// in the USB spec for details.
	PollInterval time.Duration
	// IsoSyncType is the isochronous endpoint synchronization type, as defined by USB spec.
	IsoSyncType IsoSyncType
	// UsageType is the isochronous or interrupt endpoint usage type, as defined by USB spec.
	UsageType UsageType
}

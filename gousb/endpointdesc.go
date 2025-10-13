package gousb

import (
	"fmt"
	"strings"
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

func (ep libusbEndpointDescriptor) endpointDesc(dev *DeviceDesc) EndpointDesc {
	ei := EndpointDesc{
		Address:       EndpointAddress(ep.bEndpointAddress),
		Number:        int(ep.bEndpointAddress & endpointNumMask),
		Direction:     EndpointDirection((ep.bEndpointAddress & endpointDirectionMask) != 0),
		TransferType:  TransferType(ep.bmAttributes & transferTypeMask),
		MaxPacketSize: int(ep.wMaxPacketSize),
	}
	if ei.TransferType == TransferTypeIsochronous {
		// bits 0-10 identify the packet size, bits 11-12 are the number of additional transactions per microframe.
		// Don't use libusb_get_max_iso_packet_size, as it has a bug where it returns the same value
		// regardless of alternative setting used, where different alternative settings might define different
		// max packet sizes.
		// See http://libusb.org/ticket/77 for more background.
		ei.MaxPacketSize = int(ep.wMaxPacketSize) & 0x07ff * (int(ep.wMaxPacketSize)>>11&3 + 1)
		ei.IsoSyncType = IsoSyncType(ep.bmAttributes & isoSyncTypeMask)
		// switch ep.bmAttributes & usageTypeMask {
		// case C.LIBUSB_ISO_USAGE_TYPE_DATA:
		// 	ei.UsageType = IsoUsageTypeData
		// case C.LIBUSB_ISO_USAGE_TYPE_FEEDBACK:
		// 	ei.UsageType = IsoUsageTypeFeedback
		// case C.LIBUSB_ISO_USAGE_TYPE_IMPLICIT:
		// 	ei.UsageType = IsoUsageTypeImplicit
		// }
		ei.UsageType = UsageType(ep.bmAttributes & usageTypeMask)
	}
	switch {
	// If the device conforms to USB1.x:
	//   Interval for polling endpoint for data transfers. Expressed in
	//   milliseconds.
	//   This field is ignored for bulk and control endpoints. For
	//   isochronous endpoints this field must be set to 1. For interrupt
	//   endpoints, this field may range from 1 to 255.
	// Note: in low-speed mode, isochronous transfers are not supported.
	case dev.Spec < VersionToBCD(2, 0):
		ei.PollInterval = time.Duration(ep.bInterval) * time.Millisecond

	// If the device conforms to USB[23].x and the device is in low or full
	// speed mode:
	//   Interval for polling endpoint for data transfers.  Expressed in
	//   frames (1ms)
	//   For full-speed isochronous endpoints, the value of this field should
	//   be 1.
	//   For full-/low-speed interrupt endpoints, the value of this field may
	//   be from 1 to 255.
	// Note: in low-speed mode, isochronous transfers are not supported.
	case dev.Speed == SpeedUnknown || dev.Speed == SpeedLow || dev.Speed == SpeedFull:
		ei.PollInterval = time.Duration(ep.bInterval) * time.Millisecond

	// If the device conforms to USB[23].x and the device is in high speed
	// mode:
	//   Interval is expressed in microframe units (125 µs).
	//   For high-speed bulk/control OUT endpoints, the bInterval must
	//   specify the maximum NAK rate of the endpoint. A value of 0 indicates
	//   the endpoint never NAKs. Other values indicate at most 1 NAK each
	//   bInterval number of microframes. This value must be in the range
	//   from 0 to 255.
	case dev.Speed == SpeedHigh && ei.TransferType == TransferTypeBulk:
		ei.PollInterval = time.Duration(ep.bInterval) * 125 * time.Microsecond

	// If the device conforms to USB[23].x and the device is in high speed
	// mode:
	//   For high-speed isochronous endpoints, this value must be in
	//   the range from 1 to 16. The bInterval value is used as the exponent
	//   for a 2bInterval-1 value; e.g., a bInterval of 4 means a period
	//   of 8 (2^(4-1)).
	//   For high-speed interrupt endpoints, the bInterval value is used as
	//   the exponent for a 2bInterval-1 value; e.g., a bInterval of 4 means
	//   a period of 8 (2^(4-1)). This value must be from 1 to 16.
	// If the device conforms to USB3.x and the device is in SuperSpeed mode:
	//   Interval for servicing the endpoint for data transfers. Expressed in
	//   125-µs units.
	//   For Enhanced SuperSpeed isochronous and interrupt endpoints, this
	//   value shall be in the range from 1 to 16. However, the valid ranges
	//   are 8 to 16 for Notification type Interrupt endpoints. The bInterval
	//   value is used as the exponent for a 2(^bInterval-1) value; e.g., a
	//   bInterval of 4 means a period of 8 (2^(4-1) → 2^3 → 8).
	//   This field is reserved and shall not be used for Enhanced SuperSpeed
	//   bulk or control endpoints.
	case dev.Speed == SpeedHigh || dev.Speed == SpeedSuper:
		ei.PollInterval = 125 * time.Microsecond << (ep.bInterval - 1)
	}
	return ei
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

// String returns the human-readable description of the endpoint.
func (e EndpointDesc) String() string {
	ret := make([]string, 0, 3)
	ret = append(ret, fmt.Sprintf("ep #%d %s (address %s) %s", e.Number, e.Direction, e.Address, e.TransferType))
	switch e.TransferType {
	case TransferTypeIsochronous:
		ret = append(ret, fmt.Sprintf("- %s %s", e.IsoSyncType, e.UsageType))
	case TransferTypeInterrupt:
		ret = append(ret, fmt.Sprintf("- %s", e.UsageType))
	}
	ret = append(ret, fmt.Sprintf("[%d bytes]", e.MaxPacketSize))
	return strings.Join(ret, " ")
}

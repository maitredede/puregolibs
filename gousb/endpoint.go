package gousb

import "time"

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
	case dev.Spec < Version(2, 0):
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

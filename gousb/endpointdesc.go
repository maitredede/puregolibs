package gousb

import "time"

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

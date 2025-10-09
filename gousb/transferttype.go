package gousb

// TransferType defines the endpoint transfer type.
type TransferType uint8

// Transfer types defined by the USB spec.
const (
	TransferTypeControl     TransferType = 0 // Control transfer
	TransferTypeIsochronous TransferType = 1 // Isochronous transfer
	TransferTypeBulk        TransferType = 2 // Bulk transfer
	TransferTypeInterrupt   TransferType = 3 // Interrupt transfer
	TransferTypeBulkStream  TransferType = 4 // Bulk stream transfer
	transferTypeMask                     = 0x03
)

var transferTypeDescription = map[TransferType]string{
	TransferTypeControl:     "control",
	TransferTypeIsochronous: "isochronous",
	TransferTypeBulk:        "bulk",
	TransferTypeInterrupt:   "interrupt",
}

// String returns a human-readable name of the endpoint transfer type.
func (tt TransferType) String() string {
	return transferTypeDescription[tt]
}

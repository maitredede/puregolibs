package gousb

// IsoSyncType defines the isochronous transfer synchronization type.
type IsoSyncType uint8

// Synchronization types defined by the USB spec.
const (
	IsoSyncTypeNone     IsoSyncType = 0x0 << 2
	IsoSyncTypeAsync    IsoSyncType = 0x1 << 2
	IsoSyncTypeAdaptive IsoSyncType = 0x2 << 2
	IsoSyncTypeSync     IsoSyncType = 0x3 << 2
	isoSyncTypeMask                 = 0x0C
)

var isoSyncTypeDescription = map[IsoSyncType]string{
	IsoSyncTypeNone:     "unsynchronized",
	IsoSyncTypeAsync:    "asynchronous",
	IsoSyncTypeAdaptive: "adaptive",
	IsoSyncTypeSync:     "synchronous",
}

// String returns a human-readable description of the synchronization type.
func (ist IsoSyncType) String() string {
	return isoSyncTypeDescription[ist]
}

package gousb

// UsageType defines the transfer usage type for isochronous and interrupt
// transfers.
type UsageType uint8

// Usage types for iso and interrupt transfers, defined by the USB spec.
const (
	// Note: USB3.0 defines usage type for both isochronous and interrupt
	// endpoints, with the same constants representing different usage types.
	// UsageType constants do not correspond to bmAttribute values.
	UsageTypeUndefined UsageType = iota
	IsoUsageTypeData
	IsoUsageTypeFeedback
	IsoUsageTypeImplicit
	InterruptUsageTypePeriodic
	InterruptUsageTypeNotification
	usageTypeMask = 0x30
)

var usageTypeDescription = map[UsageType]string{
	UsageTypeUndefined:             "undefined usage",
	IsoUsageTypeData:               "data",
	IsoUsageTypeFeedback:           "feedback",
	IsoUsageTypeImplicit:           "implicit data",
	InterruptUsageTypePeriodic:     "periodic",
	InterruptUsageTypeNotification: "notification",
}

func (ut UsageType) String() string {
	return usageTypeDescription[ut]
}

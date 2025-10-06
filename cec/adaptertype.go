package cec

type AdapterType int32

const (
	AdapterTypeUnknown         AdapterType = 0
	AdapterTypeP8External      AdapterType = 0x1
	AdapterTypeP8DaughterBoard AdapterType = 0x2
	AdapterTypeRPI             AdapterType = 0x100
	AdapterTypeTDA995x         AdapterType = 0x200
	AdapterTypeExynos          AdapterType = 0x300
	AdapterTypeLinux           AdapterType = 0x400
	AdapterTypeAOCEC           AdapterType = 0x500
	AdapterTypeIMX             AdapterType = 0x600
)

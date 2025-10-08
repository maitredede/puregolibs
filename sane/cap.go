package sane

type Cap int32

const (
	CapSoftSelect Cap = 1 << iota
	CapHardSelect
	CapSoftDetect
	CapEmulated
	CapAutomatic
	CapInactive
	CapAdvanced
)

func (c Cap) IsActive() bool {
	val := int32(c) & int32(CapInactive)
	return val == 0
}

func (c Cap) IsSettable() bool {
	val := int32(c) & int32(CapSoftSelect)
	return val != 0
}

func (c Cap) IsDetectable() bool {
	val := int32(c) & int32(CapSoftDetect)
	return val != 0
}

func (c Cap) IsAutomatic() bool {
	val := int32(c) & int32(CapAutomatic)
	return val != 0
}

func (c Cap) IsEmulated() bool {
	val := int32(c) & int32(CapEmulated)
	return val != 0
}

func (c Cap) IsAdvanced() bool {
	val := int32(c) & int32(CapAdvanced)
	return val != 0
}

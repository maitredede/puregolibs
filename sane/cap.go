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

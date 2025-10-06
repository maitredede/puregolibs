package sane

type Unit int32

const (
	UnitNone Unit = iota
	UnitPixel
	UnitBit
	UnitMm
	UnitDpi
	UnitPercent
	UnitUsec
)

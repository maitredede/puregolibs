package sane

type Unit int16

const (
	UnitNone Unit = iota
	UnitPixel
	UnitBit
	UnitMm
	UnitDpi
	UnitPercent
	UnitUsec
)

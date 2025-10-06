package sane

type SANE_Action int32

const (
	SANEActionGetValue SANE_Action = iota
	SANEActionSetValue
	SANEActionSetAuto
)

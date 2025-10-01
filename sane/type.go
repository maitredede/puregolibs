package sane

type Type int16

const (
	TypeBool Type = iota
	TypeInt
	TypeFloat
	TypeString
	TypeButton
	typeGroup // internal use only
)

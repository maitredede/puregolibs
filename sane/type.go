package sane

import "fmt"

type Type int32

const (
	TypeBool Type = iota
	TypeInt
	TypeFloat
	TypeString
	TypeButton
	TypeGroup
)

func (t Type) String() string {
	switch t {
	case TypeBool:
		return "bool"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeString:
		return "string"
	case TypeButton:
		return "button"
	case TypeGroup:
		return "group"
	}
	return fmt.Sprintf("?%d", int32(t))
}

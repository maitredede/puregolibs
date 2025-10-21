package internal

import "unsafe"

type SANE_Option_Descriptor struct {
	Name           unsafe.Pointer
	Title          unsafe.Pointer
	Desc           unsafe.Pointer
	Type           SANE_Value_Type
	Unit           SANE_Unit
	Size           SANE_Int
	Cap            SANE_Int
	ConstraintType SANE_Constraint_Type
	Constraint     unsafe.Pointer
}

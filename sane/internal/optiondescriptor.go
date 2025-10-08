package internal

type SANE_Option_Descriptor struct {
	Name           uintptr
	Title          uintptr
	Desc           uintptr
	Type           SANE_Value_Type
	Unit           SANE_Unit
	Size           SANE_Int
	Cap            SANE_Int
	ConstraintType SANE_Constraint_Type
	Constraint     uintptr
}

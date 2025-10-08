package internal

type SANE_Word = int32
type SANE_Handle = uintptr
type SANE_Byte = byte
type SANE_Int = SANE_Word
type SANE_Bool int32

func (b SANE_Bool) Go() bool {
	return b != SANE_FALSE
}

const (
	SANE_FALSE SANE_Bool = 0
	SANE_TRUE  SANE_Bool = 1
)

type SANE_Value_Type int32
type SANE_Unit int32
type SANE_Constraint_Type int32

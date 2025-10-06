package sane

type ConstraintType int32

const (
	ConstraintNone ConstraintType = iota
	ConstraintRange
	ConstraintWordList
	ConstraintStringList
)

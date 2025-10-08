package sane

import "fmt"

type ConstraintType int32

const (
	ConstraintNone ConstraintType = iota
	ConstraintRange
	ConstraintWordList
	ConstraintStringList
)

func (c ConstraintType) String() string {
	switch c {
	case ConstraintNone:
		return ""
	case ConstraintRange:
		return "range"
	case ConstraintWordList:
		return "numberList"
	case ConstraintStringList:
		return "stringList"
	}
	return fmt.Sprintf("?%d", c)
}

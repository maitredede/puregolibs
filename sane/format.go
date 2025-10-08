package sane

import "fmt"

type Format int32

const (
	FrameGray Format = iota
	FrameRGB
	FrameRed
	FrameGreen
	FrameBlue
)

func (f Format) String() string {
	switch f {
	case FrameGray:
		return "gray"
	case FrameRGB:
		return "rgb"
	case FrameRed:
		return "red"
	case FrameGreen:
		return "green"
	case FrameBlue:
		return "blue"
	}
	return fmt.Sprintf("?%d", int32(f))
}

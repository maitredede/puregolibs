package sane

import "fmt"

type Frame int32

const (
	FrameGray Frame = iota
	FrameRGB
	FrameRed
	FrameGreen
	FrameBlue
)

func (f Frame) String() string {
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

package plutobook

type ImageFormat int16

const (
	ImageFormatInvalid ImageFormat = -1
	ImageFormatARGB32  ImageFormat = 0
	ImageFormatRGB24   ImageFormat = 1
	ImageFormatA8      ImageFormat = 2
	ImageFormatA1      ImageFormat = 3
)

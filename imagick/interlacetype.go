package imagick

type InterlaceType int32

const (
	UndefinedInterlace InterlaceType = iota
	NoInterlace
	LineInterlace
	PlaneInterlace
	PartitionInterlace
	GIFInterlace
	JPEGInterlace
	PNGInterlace
)

package evdi

type evdiBuffer struct {
	id     int32
	buffer *byte
	width  int32
	height int32
	stride int32

	rects     *evdiRect
	rectCount int32
}

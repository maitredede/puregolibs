package evdi

type CursorSet struct {
	HotX         int32
	HotY         int32
	Width        uint32
	Height       uint32
	Enabled      uint8
	BufferLength uint32
	Buffer       *uint32
	PixelFormat  uint32
	Stride       uint32
}

type CursorMove struct {
	X int32
	Y int32
}

package plutobook

type Canvas struct{}

func (c *Canvas) Close() error {
	panic("TODO")
}

func NewImageCanvas(width, height int, imageFormat ImageFormat) (*Canvas, error) {
	panic("TODO")
}

func (c *Canvas) ClearSurface(r, g, b, a float32) error {
	panic("TODO")
}

func (c *Canvas) WriteToPNG(file string) error {
	panic("TODO")
}

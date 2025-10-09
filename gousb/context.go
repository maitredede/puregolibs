package gousb

type Context struct {
}

func NewContext() *Context {
	panic("skeleton")
}

func (c *Context) Close() error {
	panic("skeleton")
}

func (c *Context) Debug(level int) {
	panic("skeleton")
}

func (c *Context) OpenDevices(func(desc *DeviceDesc) bool) ([]*Device, error) {
	panic("skeleton")
}

package gousb

// ContextOptions holds parameters for Context initialization.
type ContextOptions struct {
	DeviceDiscovery DeviceDiscovery
}

// New creates a Context, taking into account the optional flags contained in ContextOptions
func (o ContextOptions) New() *Context {
	return newContextWithImpl(libusbImpl{
		discovery: o.DeviceDiscovery,
	})
}

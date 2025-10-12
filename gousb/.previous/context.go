package gousb

import (
	"fmt"
	"sync"
	"unsafe"
)

type Context struct {
	ctx      unsafe.Pointer
	ctxValid bool
	done     chan struct{}
	libusb   libusbIntf

	mu      sync.Mutex
	devices map[*Device]bool
}

// NewContext returns a new Context instance with default ContextOptions.
func NewContext() *Context {
	return ContextOptions{}.New()
}

func newContextWithImpl(impl libusbIntf) *Context {
	c, err := impl.init()
	if err != nil {
		panic(err)
	}
	ctx := &Context{
		ctx:     c,
		done:    make(chan struct{}),
		libusb:  impl,
		devices: make(map[*Device]bool),
	}
	go impl.handleEvents(ctx.ctx, ctx.done)
	return ctx
}

// Close releases the Context and all associated resources.
func (c *Context) Close() error {
	if c.ctx == nil {
		return nil
	}
	if err := c.checkOpenDevs(); err != nil {
		return err
	}
	c.done <- struct{}{}
	err := c.libusb.exit(c.ctx)
	c.ctx = nil
	return err
}

// Debug changes the debug level. Level 0 means no debug, higher levels
// will print out more debugging information.
// TODO(sebek): in the next major release, replace int levels with
// Go-typed constants.
func (c *Context) Debug(level int) {
	c.libusb.setDebug(c.ctx, level)
}

func (c *Context) closeDev(d *Device) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.libusb.close(d.handle)
	delete(c.devices, d)
}

func (c *Context) checkOpenDevs() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if l := len(c.devices); l > 0 {
		return fmt.Errorf("Context.Close called while %d Devices are still open, Close may be called only after all previously opened devices were successfuly closed", l)
	}
	return nil
}

// OpenDevices calls opener with each enumerated device.
// If the opener returns true, the device is opened and a Device is returned if the operation succeeds.
// Every Device returned (whether an error is also returned or not) must be closed.
// If there are any errors enumerating the devices,
// the final one is returned along with any successfully opened devices.
func (c *Context) OpenDevices(func(desc *DeviceDesc) bool) ([]*Device, error) {
	panic("skeleton")
}

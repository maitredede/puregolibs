package gousb

import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

type libusbContext unsafe.Pointer

type Context struct {
	ptr      libusbContext
	ptrValid bool

	mu      sync.Mutex
	devices map[*Device]bool
}

func Init() (*Context, error) {
	libInit()

	var ptr libusbContext
	ret := libusbInit(&ptr)
	err := errorFromRet(ret)
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		ptr:      ptr,
		ptrValid: true,
		devices:  make(map[*Device]bool),
	}
	return ctx, nil
}

func (c *Context) Close() error {
	libInit()

	if !c.ptrValid {
		return ErrInvalidContext
	}
	if err := c.checkOpenDevs(); err != nil {
		return err
	}
	libusbExit(c.ptr)
	c.ptrValid = false
	return nil
}

func (c *Context) checkOpenDevs() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if l := len(c.devices); l > 0 {
		return fmt.Errorf("Context.Close called while %d Devices are still open, Close may be called only after all previously opened devices were successfuly closed", l)
	}
	return nil
}

func (c *Context) SetDebug(level LogLevel) {
	libInit()

	if !c.ptrValid {
		return
	}
	libusbSetDebug(c.ptr, level)
}

func (c *Context) SetLocale(locale string) error {
	libInit()

	if !c.ptrValid {
		return ErrInvalidContext
	}
	// cLocale := append([]byte(locale), 0)
	// ret := libusbSetLocale(c.ptr, unsafe.Pointer(&cLocale[0]))
	ret := libusbSetLocale(c.ptr, locale)
	err := errorFromRet(ret)
	return err
}

// OpenDevices calls opener with each enumerated device.
// If the opener returns true, the device is opened and a Device is returned if the operation succeeds.
// Every Device returned (whether an error is also returned or not) must be closed.
// If there are any errors enumerating the devices,
// the final one is returned along with any successfully opened devices.
func (c *Context) OpenDevices(opener func(desc *DeviceDesc) bool) ([]*Device, error) {
	libInit()

	if !c.ptrValid {
		return nil, ErrInvalidContext
	}

	var lst *libusbDevice
	size := libusbGetDeviceList(c.ptr, &lst)
	if size < 0 {
		err := errorFromRet(size)
		return nil, err
	}
	defer libusbFreeDeviceList(lst, 0)

	devs := unsafe.Slice(lst, size)

	var reterr error
	var retDevices []*Device
	for _, dev := range devs {
		desc, err := c.getDeviceDesc(dev)
		// defer c.libusb.dereference(dev)
		defer libusbUnrefDevice(dev)
		if err != nil {
			reterr = err
			continue
		}

		if !opener(desc) {
			continue
		}

		//handle, err := c.libusb.open(dev)
		// if err != nil {
		// 	reterr = err
		// 	continue
		// }
		var handle libusbDeviceHandle
		ret := libusbOpen(dev, &handle)
		if err := errorFromRet(ret); err != nil {
			reterr = err
			continue
		}
		o := &Device{handle: handle, ctx: c, Desc: desc}
		retDevices = append(retDevices, o)
		c.mu.Lock()
		c.devices[o] = true
		c.mu.Unlock()

	}
	return retDevices, reterr
}

func (c *Context) getDeviceDesc(d libusbDevice) (*DeviceDesc, error) {
	var desc libusbDeviceDescriptor
	ret := libusbGetDeviceDescriptor(d, &desc)
	err := errorFromRet(ret)
	if err != nil {
		return nil, err
	}

	pathData := make([]byte, 16)
	pathLen := libusbGetPortNumbers(d, (*uint8)(unsafe.Pointer(&pathData[0])), int32(len(pathData)))

	var path []int
	var port int
	for _, nPort := range pathData[:pathLen] {
		port = int(nPort)
		path = append(path, port)
	}

	dev := &DeviceDesc{
		Bus:                  int(libusbGetBusNumber(d)),
		Address:              int(libusbGetDeviceAddress(d)),
		Port:                 port,
		Path:                 path,
		Speed:                Speed(libusbGetDeviceSpeed(d)),
		Spec:                 BCD(desc.bcdUSB),
		Device:               BCD(desc.bcdDevice),
		Vendor:               ID(desc.idVendor),
		Product:              ID(desc.idProduct),
		Class:                Class(desc.bDeviceClass),
		SubClass:             Class(desc.bDeviceSubClass),
		Protocol:             Protocol(desc.bDeviceProtocol),
		MaxControlPacketSize: int(desc.bMaxPacketSize0),
		iManufacturer:        int(desc.iManufacturer),
		iProduct:             int(desc.iProduct),
		iSerialNumber:        int(desc.iSerialNumber),
	}

	// Enumerate configurations
	cfgs := make(map[int]ConfigDesc)
	for i := 0; i < int(desc.bNumConfigurations); i++ {
		var cfg *libusbConfigDescriptor
		ret := libusbGetConfigDescriptor(d, uint8(i), &cfg)
		err := errorFromRet(ret)
		if err != nil {
			return nil, err
		}
		c := ConfigDesc{
			Number:         int(cfg.bConfigurationValue),
			SelfPowered:    (cfg.bmAttributes & selfPoweredMask) != 0,
			RemoteWakeup:   (cfg.bmAttributes & remoteWakeupMask) != 0,
			MaxPower:       2 * Milliamperes(cfg.MaxPower),
			iConfiguration: int(cfg.iConfiguration),
		}
		// at GenX speeds MaxPower is expressed in units of 8mA, not 2mA.
		if dev.Speed == SpeedSuper {
			c.MaxPower *= 4
		}

		ifaces := unsafe.Slice(cfg.iface, cfg.bNumInterfaces)
		// // var ifaces []C.struct_libusb_interface
		// *(*reflect.SliceHeader)(unsafe.Pointer(&ifaces)) = reflect.SliceHeader{
		// 	Data: uintptr(unsafe.Pointer(cfg.iface)),
		// 	Len:  int(cfg.bNumInterfaces),
		// 	Cap:  int(cfg.bNumInterfaces),
		// }
		c.Interfaces = make([]InterfaceDesc, 0, len(ifaces))
		// a map of interface numbers to a set of alternate settings numbers
		hasIntf := make(map[int]map[int]bool)
		for _, iface := range ifaces {
			if iface.numAltsetting == 0 {
				continue
			}

			alts := unsafe.Slice(iface.altsetting, iface.numAltsetting)
			// var alts []C.struct_libusb_interface_descriptor
			// *(*reflect.SliceHeader)(unsafe.Pointer(&alts)) = reflect.SliceHeader{
			// 	Data: uintptr(unsafe.Pointer(iface.altsetting)),
			// 	Len:  int(iface.num_altsetting),
			// 	Cap:  int(iface.num_altsetting),
			// }
			descs := make([]InterfaceSetting, 0, len(alts))
			for _, alt := range alts {
				i := InterfaceSetting{
					Number:     int(alt.bInterfaceNumber),
					Alternate:  int(alt.bAlternateSetting),
					Class:      Class(alt.bInterfaceClass),
					SubClass:   Class(alt.bInterfaceSubClass),
					Protocol:   Protocol(alt.bInterfaceProtocol),
					iInterface: int(alt.iInterface),
				}

				if hasIntf[i.Number][i.Alternate] {
					log.Printf("Device on bus %d address %d offered a descriptor for config %d with two different entries with the same interface number (%d) and the same alternate setting number (%d). gousb will use only the first one.", dev.Bus, dev.Address, c.Number, i.Number, i.Alternate)
					continue
				}
				if hasIntf[i.Number] == nil {
					hasIntf[i.Number] = make(map[int]bool)
				}
				hasIntf[i.Number][i.Alternate] = true

				ends := unsafe.Slice(alt.endpoint, alt.bNumEndpoints)
				// var ends []C.struct_libusb_endpoint_descriptor
				// *(*reflect.SliceHeader)(unsafe.Pointer(&ends)) = reflect.SliceHeader{
				// 	Data: uintptr(unsafe.Pointer(alt.endpoint)),
				// 	Len:  int(alt.bNumEndpoints),
				// 	Cap:  int(alt.bNumEndpoints),
				// }
				i.Endpoints = make(map[EndpointAddress]EndpointDesc, len(ends))
				for _, end := range ends {
					// epi := libusbEndpoint(end).endpointDesc(dev)
					epi := end.endpointDesc(dev)
					i.Endpoints[epi.Address] = epi
				}
				descs = append(descs, i)
			}
			c.Interfaces = append(c.Interfaces, InterfaceDesc{
				Number:      descs[0].Number,
				AltSettings: descs,
			})
		}
		libusbFreeConfigDescriptor(cfg)
		cfgs[c.Number] = c
	}

	dev.Configs = cfgs
	return dev, nil
}

func (c *Context) closeDev(d *Device) {
	c.mu.Lock()
	defer c.mu.Unlock()
	libusbClose(d.handle)
	delete(c.devices, d)
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
	"unsafe"

	"github.com/maitredede/puregolibs/gousb"
	"github.com/maitredede/puregolibs/gousb/usbid"
	"github.com/maitredede/puregolibs/libudev"
	"github.com/maitredede/puregolibs/tools"
	"golang.org/x/sys/unix"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	u := libudev.New()
	defer libudev.Unref(u)

	libudev.SetLogPriority(u, libudev.LogDebug)
	libudev.SetLogFn(u, func(udev libudev.UDev, priority int32, file string, line int32, fn, format string, args unsafe.Pointer) {
		fmt.Printf("level=%v file=%v line=%v fn=%v format=%v args=%v\n", priority, file, line, fn, format, args)
	})

	m := libudev.MonitorNewFromNetlink(u, "udev")
	fmt.Printf("monitor m=%v\n", m)
	defer libudev.MonitorUnref(m)

	ret := libudev.MonitorFilterAddMatchSubsystemDevType(m, "usb", "")
	fmt.Printf("addMatch ret=%v\n", ret)
	ret = libudev.MonitorEnableReceiving(m)
	fmt.Printf("enableReceiving ret=%v\n", ret)

	fd := int(libudev.MonitorGetFd(m))
	timeout := tools.DurationToTimeVal(1 * time.Second)

	func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			fds := unix.FdSet{}
			fds.Zero()
			fds.Set(fd)

			n, err := unix.Select(fd+1, &fds, nil, nil, &timeout)
			if err != nil {
				fmt.Printf("select error: %v\n", err)
				continue
			}
			if n <= 0 {
				continue
			}
			if fds.IsSet(fd) {
				dev := libudev.MonitorReceiveDevice(m)
				processDevice(dev)
			}
		}
	}()
}

func processDevice(dev libudev.Device) {
	defer libudev.DeviceUnref(dev)

	node := libudev.DeviceGetDevNode(dev)
	if len(node) > 0 {

		action := libudev.DeviceGetAction(dev)
		if len(action) == 0 {
			action = "exists"
		}
		vendor := libudev.DeviceGetSysAttrValue(dev, "idVendor")
		if len(vendor) == 0 {
			vendor = "0000"
		}
		product := libudev.DeviceGetSysAttrValue(dev, "idProduct")
		if len(product) == 0 {
			product = "0000"
		}

		var vendorString string
		var productString string

		vendorID, err := strconv.ParseInt(vendor, 16, 32)
		var vendorEntry *usbid.Vendor
		if err == nil {
			var ok bool
			vendorEntry, ok = usbid.Vendors[gousb.ID(vendorID)]
			if ok {
				vendorString = vendorEntry.Name
			}
		}
		productID, err := strconv.ParseInt(product, 16, 32)
		if err == nil && vendorEntry != nil {
			productEntry, ok := vendorEntry.Product[gousb.ID(productID)]
			if ok {
				productString = productEntry.Name
			}
		}

		fmt.Printf("%s %s %6s %s:%s %s (%s - %s)\n",
			libudev.DeviceGetSubsystem(dev),
			libudev.DeviceGetDevType(dev),
			action,
			vendor,
			product,
			node,
			vendorString,
			productString,
		)
	}
}

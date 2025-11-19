package libudev

import (
	"strconv"
	"testing"
	"unsafe"

	"github.com/maitredede/puregolibs/gousb"
	"github.com/maitredede/puregolibs/gousb/usbid"
	"github.com/stretchr/testify/assert"
)

func processDevice(t *testing.T, dev Device) {
	defer DeviceUnref(dev)

	node := DeviceGetDevNode(dev)
	if len(node) > 0 {

		action := DeviceGetAction(dev)
		if len(action) == 0 {
			action = "exists"
		}
		vendor := DeviceGetSysAttrValue(dev, "idVendor")
		if len(vendor) == 0 {
			vendor = "0000"
		}
		product := DeviceGetSysAttrValue(dev, "idProduct")
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

		// t.Logf("found device path=%s", path)
		// t.Logf(" node: %s", node)
		t.Logf("%s %s %6s %s:%s %s (%s - %s)",
			DeviceGetSubsystem(dev),
			DeviceGetDevType(dev),
			action,
			vendor,
			product,
			node,
			vendorString,
			productString,
		)
	}
}

func TestEnumerateDevices(t *testing.T) {
	u := New()
	defer Unref(u)
	SetLogPriority(u, LogDebug)
	SetLogFn(u, func(udev UDev, priority int32, file string, line int32, fn, format string, args unsafe.Pointer) {
		t.Logf("level=%v file=%v line=%v fn=%v format=%v args=%v", priority, file, line, fn, format, args)
	})

	e := EnumerateNew(u)
	assert.NotNil(t, e, "enumerate should be allocated")
	defer EnumerateUnref(e)

	EnumerateAddMatchSubsystem(e, "usb")

	n := EnumerateScanDevices(e)
	t.Logf("scan returned %v", n)

	devices := EnumerateGetListEntry(e)

	ListEntryForEach(devices, func(entry ListEntry) {
		path := ListEntryGetName(entry)

		dev := DeviceNewFromSyspath(u, path)
		processDevice(t, dev)
	})
}

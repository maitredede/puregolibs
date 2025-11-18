package libudev

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

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

		t.Logf("found device path=%s", path)

		// dev := DeviceNewFromSyspath(u, path)
	})
}

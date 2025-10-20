package libnfc

import (
	"testing"

	"github.com/maitredede/puregolibs/strings"
)

func TestCustomDriver(t *testing.T) {
	t.Setenv("LIBNFC_LOG_LEVEL", "3")
	t.Setenv("LIBNFC_AUTO_SCAN", "true")
	t.Setenv("LIBNFC_INTRUSIVE_SCAN", "true")
	testDriver := &Driver{
		Name:     "go_test",
		ScanMode: ScanModeNotIntrusive,
	}
	if err := RegisterDriver(testDriver); err != nil {
		t.Fatal(err)
	}

	nfc, err := InitContext()
	if err != nil {
		t.Fatal(err)
	}
	defer nfc.Close()

	devices, err := nfc.ListDevices()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d devices: %v", len(devices), devices)
}

func TestListDrivers(t *testing.T) {
	t.Setenv("LIBNFC_LOG_LEVEL", "3")
	t.Setenv("LIBNFC_AUTO_SCAN", "true")
	t.Setenv("LIBNFC_INTRUSIVE_SCAN", "true")
	libInit()

	nfc, err := InitContext()
	if err != nil {
		t.Fatal(err)
	}
	defer nfc.Close()

	testDriver := &Driver{
		Name:     "go_test",
		ScanMode: ScanModeNotIntrusive,
	}
	if err := RegisterDriver(testDriver); err != nil {
		t.Fatal(err)
	}

	drivers, err := GetNativeDrivers()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d native drivers:", len(drivers))
	for _, nd := range drivers {
		name := strings.GoStringN(uintptr(nd.name), 256)
		t.Logf("driver: %s (scan %s)", name, nd.scanType)
	}
}

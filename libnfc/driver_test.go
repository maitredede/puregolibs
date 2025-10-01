package libnfc

import (
	"os"
	"testing"
)

func TestCustomDriver(t *testing.T) {
	os.Setenv("LIBNFC_LOG_LEVEL", "3")
	os.Setenv("LIBNFC_AUTO_SCAN", "true")
	os.Setenv("LIBNFC_INTRUSIVE_SCAN", "true")
	testDriver := &Driver{
		Name: "test",
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

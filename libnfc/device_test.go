package libnfc

import "testing"

func TestDeviceOpenClose(t *testing.T) {
	nfc, err := InitContext()
	if err != nil {
		t.Fatal(err)
	}
	defer nfc.Close()

	dev, err := nfc.OpenDefault()
	if err != nil {
		t.Fatal(err)
	}
	name := dev.Name()
	t.Logf("opened nfc device: %s", name)
	connString, err := dev.ConnString()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("connString: %s", connString)

	infos, err := dev.GetInformationAbout()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("infos: %s", infos)

	if err := dev.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDeviceList(t *testing.T) {
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

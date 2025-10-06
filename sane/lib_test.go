package sane

import "testing"

func TestInitExit(t *testing.T) {

	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	Exit()
}

func TestGetDevices(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Exit)

	devices, err := GetDevices(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d devices: %v", len(devices), devices)
}

func TestOpenClose(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Exit)

	devices, err := GetDevices(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d devices: %v", len(devices), devices)
	if len(devices) == 0 {
		t.Fatal("no devices found")
	}
	first := devices[0]
	t.Logf("using device %+v", first)

	h, err := Open(first.Name)
	if err != nil {
		t.Fatal(err)
	}
	Close(h)
}

func TestGetOptionDescriptors(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(Exit)

	devices, err := GetDevices(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d devices: %v", len(devices), devices)
	if len(devices) == 0 {
		t.Fatal("no devices found")
	}
	first := devices[0]
	t.Logf("using device %+v", first)

	h, err := Open(first.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer Close(h)

	options, err := GetOptionDescriptors(h)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("options: %+v", options)
}

package sane

import (
	"fmt"
	"testing"
)

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
	h.Close()
}

func TestGetOptionDescriptorsParams(t *testing.T) {
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
	defer h.Close()

	descList, err := h.GetOptionDescriptors()
	if err != nil {
		t.Fatal(err)
	}
	var curGroup string
	for _, desc := range descList {
		if desc.Type == TypeGroup {
			curGroup = desc.Title
			continue
		}

		s := fmt.Sprintf("%s\t%2d %s %s %d", curGroup, desc.Number, desc.Title, desc.Type, desc.BinSize)

		value, err := desc.GetValue()
		if err != nil {
			t.Log(err)
			t.Fail()
			continue
		}
		s = fmt.Sprintf("%s\tv=%v", s, value)
		if desc.Constraint != nil {
			s = fmt.Sprintf("%s\tconstraint %s=%v", s, desc.ConstraintType, desc.Constraint)
		}

		t.Log(s)
	}

	params, err := h.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", params)
}

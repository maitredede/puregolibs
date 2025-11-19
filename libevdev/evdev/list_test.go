package evdev

import (
	"testing"
)

func TestList(t *testing.T) {
	devices, err := ListDevicePaths()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("found %d devices:", len(devices))
	for _, d := range devices {
		t.Logf(" %v (%v)", d.Path, d.Name)
	}
}

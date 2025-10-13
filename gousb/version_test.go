package gousb

import "testing"

func TestVersion(t *testing.T) {
	v := GetVersion()
	t.Logf("libusb version: %+v", v)
}

package evdi

import "testing"

func TestIsXorgRunning(t *testing.T) {
	r := IsXorgRunning()
	t.Logf("xorg running: %v", r)
}

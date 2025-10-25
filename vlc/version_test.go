package vlc

import "testing"

func TestVersion(t *testing.T) {
	v := GetVersion()

	t.Logf("libvlc version: %v", v)
}

package fontconfig

import "testing"

func TestVersion(t *testing.T) {
	t.Logf("libfontconfig version: %d", GetVersion())
	t.Logf("libfontconfig version string: %s", GetVersionString())
}

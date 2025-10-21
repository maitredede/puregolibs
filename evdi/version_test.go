package evdi

import "testing"

func TestVersion(t *testing.T) {
	v := VersionString()
	t.Logf("evdi version: %s", v)
}

package libnfc

import "testing"

func TestVersion(t *testing.T) {
	v := Version()
	t.Logf("version: %s", v)
}

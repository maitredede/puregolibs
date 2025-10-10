package plutobook

import (
	"testing"

	"github.com/jupiterrider/ffi"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {

	ffiVersionString := ffi.GetVersion()
	t.Logf("using ffi version: %s", ffiVersionString)

	minVersion := 0*10000 + 9*100 + 0

	v := VersionNumber()
	t.Logf("version num: %d", v)
	assert.GreaterOrEqual(t, v, minVersion)
}

func TestVersionString(t *testing.T) {
	ffiVersionString := ffi.GetVersion()
	t.Logf("using ffi version: %s", ffiVersionString)
	t.Logf("plutobook version: %s", Version())
}

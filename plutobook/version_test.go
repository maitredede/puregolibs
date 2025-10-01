package plutobook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {

	minVersion := 0*10000 + 9*100 + 0

	v := VersionNumber()
	t.Logf("version num: %d", v)
	assert.GreaterOrEqual(t, minVersion, v)
}

func TestVersionString(t *testing.T) {
	t.Logf("version string: %s", Version())
}

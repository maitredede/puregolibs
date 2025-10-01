package plutobook

import "testing"

func TestBuildInfo(t *testing.T) {
	t.Logf("build info: %s", BuildInfo())
}

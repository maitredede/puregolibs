package fontconfig

import (
	"testing"
)

func TestConfigHome(t *testing.T) {
	t.Logf("fc config home: %s", ConfigHome())
}

package dotnetcore

import (
	"slices"
	"testing"
)

func TestRuntimeInit(t *testing.T) {
	sdks := locateSDK()
	if len(sdks) == 0 {
		t.Fatal("no sdks found")
	}

	slices.Sort(sdks)

	r, err := InitializeRuntime(sdks[len(sdks)-1])
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
}

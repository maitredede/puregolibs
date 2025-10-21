package evdi

import "testing"

func TestSetLogging(t *testing.T) {
	SetLogging(func(s string) {
		t.Log(s)
	})

	libEvdiOpen(42)
}

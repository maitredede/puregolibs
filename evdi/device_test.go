package evdi

import "testing"

func TestOpenClose(t *testing.T) {
	SetLogging(func(s string) {
		t.Log(s)
	})

	dev, err := OpenAttachedToNone()
	if err != nil {
		t.Fatal(err)
	}

	if err := dev.Close(); err != nil {
		t.Fatal(err)
	}
}

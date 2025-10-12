package gousb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextInitExit(t *testing.T) {
	ctx, err := Init()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotZero(t, ctx.ptr, "pointer should have a value")

	ctx.SetDebug(LogLevelDebug)

	err = ctx.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocales(t *testing.T) {
	t.Skip("locales test for libusb are failing")

	ctx, err := Init()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()
	ctx.SetDebug(LogLevelDebug)

	// https://github.com/libusb/libusb/blob/v1.0.28/libusb/strerror.c#L51
	locales := []string{"en", "nl", "fr", "ru", "de", "hu"}

	for _, loc := range locales {
		err := ctx.SetLocale(loc)
		t.Logf("locale: %s err=%v", loc, err)
		if err != nil {
			t.Fail()
		}
	}

}

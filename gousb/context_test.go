package gousb

import (
	"strings"
	"testing"
	"unicode"

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

func TestContextSetOption(t *testing.T) {
	ctx, err := Init()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	if err := ctx.SetOptionLogLevel(LogLevelDebug); err != nil {
		t.Fatal(err)
	}
	if err := ctx.SetOptionLogCallback(func(ctx *Context, level LogLevel, str string) {
		msg := strings.TrimRightFunc(str, unicode.IsSpace)
		t.Logf("logcb: [%s] %s", level, msg)
	}); err != nil {
		t.Fatal(err)
	}

	_, err = ctx.OpenDevices(func(desc *DeviceDesc) bool {
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestContextOpenDevices(t *testing.T) {
	ctx, err := Init()
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()

	devices, err := ctx.OpenDevices(func(desc *DeviceDesc) bool {
		return true
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Logf("opened %d devices", len(devices))
	for _, d := range devices {

		t.Logf(" vid=%s pid=%s", d.Desc.Vendor, d.Desc.Product)

		d.Close()
	}
}

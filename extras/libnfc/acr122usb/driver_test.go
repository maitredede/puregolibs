package acr122usb

import (
	gostrings "strings"
	"testing"
	"unicode"

	"github.com/maitredede/puregolibs/gousb"
	"github.com/maitredede/puregolibs/libnfc"
	"github.com/stretchr/testify/assert"
)

func TestDriverScan(t *testing.T) {
	t.Setenv("LIBNFC_LOG_LEVEL", "3")
	t.Setenv("LIBNFC_AUTO_SCAN", "true")
	t.Setenv("LIBNFC_INTRUSIVE_SCAN", "true")

	usb, err := gousb.Init(
		gousb.WithLogLevel(gousb.LogLevelDebug),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer usb.Close()
	if err := usb.SetOptionLogCallback(func(ctx *gousb.Context, level gousb.LogLevel, str string) {
		msg := gostrings.TrimRightFunc(str, unicode.IsSpace)
		t.Logf("usb: [%s] %s", level, msg)
	}); err != nil {
		t.Fatal(err)
	}

	if err := RegisterGoACR122usbDriver(usb); err != nil {
		t.Fatal(err)
	}

	nfc, err := libnfc.InitContext()
	if err != nil {
		t.Fatal(err)
	}
	defer nfc.Close()

	nativeDrivers, err := libnfc.GetNativeDrivers()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, nativeDrivers, 1, "only the golang driver should be registered")

	devices, err := nfc.ListDevices()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d devices", len(devices))
	t.Logf("%+v", devices)
}

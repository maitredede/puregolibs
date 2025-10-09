package cec

import (
	"os"
	"testing"

	"github.com/jupiterrider/ffi"
)

func TestInitExit(t *testing.T) {
	e := os.Getenv("FFI_NO_EMBED")
	t.Logf("FFI_NO_EMBED=%s", e)
	v := ffi.GetVersion()
	vn := ffi.GetVersionNumber()
	t.Logf("using libffi: %s (%d)", v, vn)

	libInit()

	callbacks := Callbacks{
		LogMessage: func(cbparam any, message LogMessage) {
			t.Logf("log: %v: %s", message.Level, message.Message)
		},
		KeyPress: func(cbparam any, key Keypress) {
			t.Logf("keypress: %v", key)
		},
		CommandReceived: func(cbparam any, command Command) {
			t.Logf("cmdReceived: %v", command)
		},
		ConfigurationChanged: func(cbparam any, configuration Configuration) {
			t.Logf("cfgChanged: %v", configuration)
		},
		Alert: func(cbparam any, alert Alert, param Parameter) {
			t.Logf("alert: %v %v", alert, param)
		},
		MenuStateChanged: func(cbparam any, state MenuState) int32 {
			t.Logf("menuStateChanged: %v", state)
			return 0
		},
		SourceActivated: func(cbparam any, logicalAddress LogicalAddress, activated bool) {
			t.Logf("sourceActivated: %v => %v", logicalAddress, activated)
		},
		CommandHandler: func(cbparam any, command Command) int32 {
			t.Logf("commandHandler: %v", command)
			return 0
		},
	}

	con, err := Open("", "", callbacks)
	if err != nil {
		t.Fatal(err)
	}

	info, err := con.GetLibInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("libinfo: %s", info)
	t.Logf("config:\n%+v", con.cfg)

	devices := con.GetActiveDevices()
	t.Logf("active devices: %+v", devices)

	if err := con.Close(); err != nil {
		t.Fatal(err)
	}
}

package cec

import "testing"

func TestInitExit(t *testing.T) {
	cfg := Configuration{
		ClientVersion: VersionCurrent,
		DeviceName:    t.Name(),
		DeviceTypes: DeviceTypeList{
			DeviceTypeTV,
		},
		AutodectAddress: true,
	}

	con, err := Initialise(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	info, err := con.GetLibInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("libinfo: %s", info)
	t.Logf("config:\n%+v", con.cfg)

	if err := con.Close(); err != nil {
		t.Fatal(err)
	}
}

package cec

import "testing"

func TestCecVersionToString(t *testing.T) {
	versions := []CecVersion{
		CecVersion_1_2,
		CecVersion_1_2A,
		CecVersion_1_3,
		CecVersion_1_3A,
		CecVersion_1_4,
		CecVersion_2_0,
	}

	for _, v := range versions {
		t.Logf("version value=%d string=%s", int(v), v.String())
	}
}

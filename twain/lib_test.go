package twain

import "testing"

func TestLib(t *testing.T) {
	ret := libDSMEntry(nil, nil, DataGroupControl, DataArgumentImageInfo, MessageGet, 0)
	_ = ret
}

package libc

import "testing"

func TestAllocFree(t *testing.T) {
	p := CAlloc(1, 1)
	if p == nil {
		t.Fatal("memory allocation failed")
	}
	Free(p)
}

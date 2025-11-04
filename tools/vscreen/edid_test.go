package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEDIDLength(t *testing.T) {
	edid := buildEdid()

	assert.Len(t, edid, 128, "EDID should have length of 128")
}

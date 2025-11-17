package libudev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDevNewUnref(t *testing.T) {
	u := New()
	assert.NotNil(t, u, "udev object should exists")
	assert.Len(t, infoTrack, 1, "tracking should contain item")
	Unref(u)
	assert.Len(t, infoTrack, 0, "tracking should not contain item")
}

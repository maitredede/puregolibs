package libevdev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvdevNewUnref(t *testing.T) {
	u := New()
	assert.NotNil(t, u, "evdev object should exists")
	Free(u)
}

package libudev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDevNewUnref(t *testing.T) {
	u := New()
	assert.NotNil(t, u, "udev object should exists")
	Unref(u)
}

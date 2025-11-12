package libevdev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitUnload(t *testing.T) {
	err := Initialize()
	assert.NoError(t, err, "init should work")
	err = Unload()
	assert.NoError(t, err, "unload should work")
}

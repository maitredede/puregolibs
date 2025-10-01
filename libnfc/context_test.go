package libnfc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextInitClose(t *testing.T) {
	v := Version()
	t.Logf("nfclib: %s", v)
	nfc, err := InitContext()
	assert.NoError(t, err)
	assert.NotZero(t, nfc.ptr)

	err = nfc.Close()
	assert.NoError(t, err)
}

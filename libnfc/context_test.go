package libnfc

import (
	"testing"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"
)

func TestContextInitClose(t *testing.T) {
	log := slogt.New(t)

	v := Version()
	t.Logf("libnfc version: %s", v)
	nfc, err := InitContext(WithSLogger(log))
	assert.NoError(t, err)
	assert.NotZero(t, nfc.ptr)

	err = nfc.Close()
	assert.NoError(t, err)
}

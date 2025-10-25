package vlc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceArguments(t *testing.T) {
	v := GetVersion()
	t.Logf("using libvlc version: %s", v)

	args := []string{
		"--intf=dummy",
		"--ignore-config",
		"--no-media-library",
		"--no-one-instance",
		"--no-osd",
		"--no-snapshot-preview",
		"--no-stats",
		"--no-video-title-show",
		"-vvv",
	}

	instance, err := New(args)
	assert.NoError(t, err, "no error should occure")
	assert.NotNil(t, instance.ptr, "instance should exist")
	err = instance.Close()
	assert.NoError(t, err, "no error should occure")
}

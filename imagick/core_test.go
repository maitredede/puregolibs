package imagick

import "testing"

func TestCore(t *testing.T) {
	coreEnv := MagickCoreGenesis()
	coreEnv.Terminus()
}

package fontconfig

import "unsafe"

type fcConfigPtr unsafe.Pointer

func ConfigHome() string {
	libInit()

	return libFcConfigHome()
}

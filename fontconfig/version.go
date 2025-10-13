package fontconfig

import "fmt"

func GetVersion() int {
	libInit()
	return int(libFcGetVersion())
}

func GetVersionString() string {
	v := GetVersion()
	major := v / 10000
	minorBase := v - major*10000
	minor := minorBase / 100
	patch := (minorBase - minor*100)
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

package dotnetcore

import "testing"

func TestListSDKVersions(t *testing.T) {
	sdkPaths := locateSDK()

	t.Logf("%+v", sdkPaths)

	// versions := make(map[string]string)
	// for _, p := range sdkPaths {

	// }
}

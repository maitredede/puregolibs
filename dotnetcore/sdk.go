package dotnetcore

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	linuxSDKPaths = []string{
		"/usr/share/dotnet/shared/Microsoft.NETCore.App",
		"$HOME/.dotnet/shared/Microsoft.NETCore.App",
	}

	darwinSDKPaths = []string{
		"/usr/local/share/dotnet/shared/Microsoft.NETCore.App",
		"$HOME/.dotnet/shared/Microsoft.NETCore.App",
	}
)

// locateSDK finds the SDK path
// TODO: allow the user to use a specific version when multiple SDKs are present.
func locateSDK() (sdkDirectories []string) {
	var basePaths []string
	switch runtime.GOOS {
	case "darwin":
		basePaths = darwinSDKPaths
		break
	case "linux":
		basePaths = linuxSDKPaths
	}
	// Replace HOME env var from base paths:
	homeEnv := os.Getenv("HOME")

	for _, basePath := range basePaths {
		basePath = strings.Replace(basePath, "$HOME", homeEnv, 1)
		directories, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}
		for _, d := range directories {
			fullPath := filepath.Join(basePath, d.Name())
			sdkDirectories = append(sdkDirectories, fullPath)
		}
	}
	return sdkDirectories
}

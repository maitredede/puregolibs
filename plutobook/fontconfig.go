package plutobook

// SetFontConfigPath Sets the `FONTCONFIG_PATH` environment variable for the current process.
func SetFontConfigPath(path string) {
	libInit()

	libSetFontconfigPath(path)
}

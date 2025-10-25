package vlc

// GetVersion Retrieve libvlc version
func GetVersion() string {
	libInit()

	return libvlcGetVersion()
}

// GetCompiler Retrieve libvlc compiler version
func GetCompiler() string {
	libInit()

	return libvlcGetCompiler()
}

// GetChangeset Retrieve libvlc changeset
func GetChangeset() string {
	libInit()

	return libvlcGetChangeset()
}

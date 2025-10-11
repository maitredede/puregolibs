package plutobook

// SetSSLCaInfo Sets the path to a file containing trusted CA certificates
func SetSSLCaInfo(path string) {
	libInit()

	libSetSSLCaInfo(path)
}

// SetSSLCaPath Sets the path to a directory containing trusted CA certificates
func SetSSLCaPath(path string) {
	libInit()

	libSetSSLCaPath(path)
}

// SetSSLVerifyPeer Enables or disables SSL peer certificate verification
func SetSSLVerifyPeer(verify bool) {
	libInit()

	libSetSSLVerifyPeer(verify)
}

// SetSSLVerifyHost Enables or disables SSL host name verification
func SetSSLVerifyHost(verify bool) {
	libInit()

	libSetSSLVerifyHost(verify)
}

func SetHttpFollowRedirects(follow bool) {
	libInit()

	libSetHttpFollowRedirects(follow)
}

func SetHttpMaxRedirects(amount int) {
	libInit()

	libSetHttpMaxRedirects(int32(amount))
}

func SetHttpTimeout(timeout int) {
	libInit()

	libSetHttpTimeout(int32(timeout))
}

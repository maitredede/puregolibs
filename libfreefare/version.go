package libfreefare

func Version() string {
	libInit()

	return libVersion()
}

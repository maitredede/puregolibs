package evdi

func IsXorgRunning() bool {
	initLib()

	return libXorgRunning()
}

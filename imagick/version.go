package imagick

func GetMagickVersion() string {
	libInit()

	return libCoreGetVersion()
}

func GetMagickVersionWand() string {
	libInit()

	return libWandGetVersion()
}

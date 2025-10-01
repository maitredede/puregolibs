package plutobook

var (
	libVersion       func() int16
	libVersionString func() string
)

func Version() string {
	libInit()
	return libVersionString()
}

func VersionNumber() int {
	libInit()
	return int(libVersion())
}

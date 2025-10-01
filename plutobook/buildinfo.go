package plutobook

var (
	libBuildInfo func() string
)

func BuildInfo() string {
	libInit()

	return libBuildInfo()
}

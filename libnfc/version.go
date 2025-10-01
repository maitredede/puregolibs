package libnfc

var (
	libVersion func() string
)

func Version() string {
	libInit()

	return libVersion()
}

package evdi

import _ "embed"

var (
	//go:embed EDIDv1_1280x800
	EDIDv1_1280x800 []byte
	//go:embed EDIDv1_1440x900
	EDIDv1_1440x900 []byte
	//go:embed EDIDv1_1600x900
	EDIDv1_1600x900 []byte
	//go:embed EDIDv1_1680x1050
	EDIDv1_1680x1050 []byte
	//go:embed EDIDv2_1280x720
	EDIDv2_1280x720 []byte
	//go:embed EDIDv2_1920x1080
	EDIDv2_1920x1080 []byte
	//go:embed EDIDv2_3840x2160
	EDIDv2_3840x2160 []byte
)

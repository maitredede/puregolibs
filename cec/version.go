package cec

import "fmt"

const VersionCurrent VersionValue = 0x060002

type VersionValue uint32

func (v VersionValue) String() string {
	major := (v & 0xFF0000) >> 16
	minor := (v & 0x00FF00) >> 8
	patch := (v & 0x0000FF) >> 0
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

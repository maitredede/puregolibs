package libnfc

import "fmt"

type ScanMode int16

const (
	ScanModeNotIntrusive ScanMode = iota
	ScanModeIntrusive
	ScanModeNotAvailable
)

func (m ScanMode) String() string {
	switch m {
	case ScanModeIntrusive:
		return "intrusive"
	case ScanModeNotIntrusive:
		return "not intrusive"
	case ScanModeNotAvailable:
		return "not available"
	}
	return fmt.Sprintf("?%d", (int16)(m))
}

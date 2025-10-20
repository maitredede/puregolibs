package libnfc

import "fmt"

type Mode byte

const (
	ModeTarget Mode = iota
	ModeInitiator
)

func (m Mode) String() string {
	switch m {
	case ModeTarget:
		return "target"
	case ModeInitiator:
		return "initiator"
	}
	return fmt.Sprintf("?%d", (byte)(m))
}

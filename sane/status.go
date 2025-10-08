package sane

type SANE_Status int32

const (
	StatusGood SANE_Status = iota
	StatusUnsupported
	StatusCancelled
	StatusDeviceBusy
	StatusInval
	StatusEOF
	StatusJammed
	StatusNoDocs
	StatusCoverOpen
	StatusIOError
	StatusNoMem
	StatusAccessDenied
)

const (
	StatusWarmingUp SANE_Status = 12
	StatusHWLocked  SANE_Status = 13
)

func (s SANE_Status) String() string {
	libInit()

	return libSaneStrStatus(s)
}

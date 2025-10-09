package cec

type Alert int32

const (
	AlertServiceDevice Alert = iota
	AlertConnectionLost
	AlertPermissionError
	AlertPortBusy
	AlertPhysicalAddressError
	AlertTvPollFailed
)

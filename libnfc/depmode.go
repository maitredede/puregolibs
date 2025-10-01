package libnfc

type DepMode int16

const (
	NDMUndefined DepMode = iota
	NDMPassive
	NDMActive
)

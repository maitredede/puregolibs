package libnfc

type BaudRate int16

const (
	NBRUndefined BaudRate = iota
	NBR106
	NBR212
	NBR424
	NBR847
)

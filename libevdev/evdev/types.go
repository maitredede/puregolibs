package evdev

// InputID ...
type InputID struct {
	BusType uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

// AbsInfo describes details on ABS input types
type AbsInfo struct {
	Value      int32
	Minimum    int32
	Maximum    int32
	Fuzz       int32
	Flat       int32
	Resolution int32
}

// InputKeymapEntry is used to retrieve and modify keymap data
type InputKeymapEntry struct {
	Flags    uint8
	Len      uint8
	Index    uint16
	KeyCode  uint32
	ScanCode [32]uint8
}

// InputMask ...
type InputMask struct {
	Type      uint32
	CodesSize uint32
	CodesPtr  uint64
}

// UinputUserDevice is used when creating or cloning a device
type UinputUserDevice struct {
	Name       [uinputMaxNameSize]byte
	ID         InputID
	EffectsMax uint32
	Absmax     [absSize]int32
	Absmin     [absSize]int32
	Absfuzz    [absSize]int32
	Absflat    [absSize]int32
}

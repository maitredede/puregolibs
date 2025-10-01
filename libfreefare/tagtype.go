package libfreefare

type TagType int16

const (
	TypeFelica TagType = iota
	TypeMifareMini
	TypeMifareClassic1K
	TypeMifareClassic4K
	TypeMifareDesfire
	// TypeMifarePlusS2K
	// TypeMifarePlusS4K
	// TypeMifarePlusX2K
	// TypeMifarePlusX4K
	TypeMifareUltralight
	TypeMifareUltralightC
	TypeNTag21x
)

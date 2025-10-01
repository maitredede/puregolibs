package libnfc

type ModulationType int16

const (
	NMTIso14443A ModulationType = iota + 1
	NMTJewel
	NMTIso14443B
	NMTIso14443BI   // pre-ISO14443B aka ISO/IEC 14443 B' or Type B'
	NMTIso14443B2SR // ISO14443-2B ST SRx
	NMTIso14443B2CT // ISO14443-2B ASK CTx
	NMTFelica
	NMTDep
	NMTBarcode          // Thinfilm NFC Barcode
	NMTISO144443BICLASS // HID iClass 14443B mode
)

package libfreefare

type MifareClassicKey [6]byte

func (k *MifareClassicKey) toArg() *byte {
	return &k[0]
}

type mad_aid struct {
	applicationCode     byte
	functionClusterCode byte
}

func (m mad_aid) toArg() *byte {
	res := [2]byte{
		m.applicationCode,
		m.functionClusterCode,
	}
	return &res[0]
}

type Mad struct {
	tag *MifareTag
	ptr uintptr
}

func (m *Mad) GetVersion() int16 {
	return libMadGetVersion(m.ptr)
}

func (m *Mad) SetVersion(version byte) {
	libMadSetVersion(m.ptr, version)
}

type MifareClassicKeyType int16

const (
	MFCKeyA MifareClassicKeyType = iota
	MFCKeyB
)

package libfreefare

var (
	// Mad		 mad_read (MifareTag tag)
	libMadRead func(tag uintptr) uintptr
	// int		 mad_get_version (Mad mad)
	libMadGetVersion func(mad uintptr) int16
	// void		 mad_set_version (Mad mad, const uint8_t version)
	libMadSetVersion func(mad uintptr, version byte)

	// ssize_t		 mifare_application_read (MifareTag tag, Mad mad, const MadAid aid, void *buf, size_t nbytes, const MifareClassicKey key, const MifareClassicKeyType key_type);
	libMifareApplicationRead func(tag uintptr, mad uintptr, aid *byte, buf uintptr, nbytes uint32, key *byte, keyType MifareClassicKeyType) int32
)

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

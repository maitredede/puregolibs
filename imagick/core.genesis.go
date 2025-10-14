package imagick

type MagickCoreEnvironment struct {
}

func MagickCoreGenesis() *MagickCoreEnvironment {
	libInit()
	libCoreGenesis()
	e := &MagickCoreEnvironment{}
	return e
}

func (e *MagickCoreEnvironment) Terminus() {
	libInit()
	libCoreTerminus()
}

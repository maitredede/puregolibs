package imagick

type MagickWandEnvironment struct {
}

func MagickWandGenesis() *MagickWandEnvironment {
	libInit()
	libWandGenesis()
	e := &MagickWandEnvironment{}
	return e
}

func (e *MagickWandEnvironment) Terminus() {
	libInit()
	libWandTerminus()
}

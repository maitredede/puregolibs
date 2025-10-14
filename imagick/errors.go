package imagick

type WandEnvError struct {
	m string
}

func (e *WandEnvError) Error() string {
	return e.m
}

var ErrWandEnvinment error = &WandEnvError{m: "missing wand environment. you need to call 'env := MagickWandGenesis()' at the beginning of your program, and 'env.Terminus()' at the end"}

type InvalidWandError struct {
	m string
}

func (e *InvalidWandError) Error() string {
	return e.m
}

var ErrInvalidWand error = &InvalidWandError{m: "invalid wand"}

type MagickException struct {
	m string
	t ExceptionType
}

func (e *MagickException) Error() string {
	return e.m
}

func (e *MagickException) Type() ExceptionType {
	return e.t
}

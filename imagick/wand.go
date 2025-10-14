package imagick

import (
	"unsafe"
)

type magickWandPtr unsafe.Pointer

type MagickWand struct {
	ptr magickWandPtr
}

func (e *MagickWandEnvironment) NewMagickWand() (*MagickWand, error) {
	libInit()
	// isEnvReady := libWandIsMagickWandInstantiated()
	// if !isEnvReady {
	// 	return nil, ErrWandEnvinment
	// }
	ptr := libWandNewMagickWand()

	if !libWandIsMagickWand(ptr) {
		return nil, ErrInvalidWand
	}

	w := &MagickWand{
		ptr: ptr,
	}
	return w, nil
}

func (w *MagickWand) Close() error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native DestroyMagickWand
	libWandDestroyMagickWand(w.ptr)
	return nil
}

func (w *MagickWand) ReadImageBlob(blob []byte) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native MagickReadImageBlob
	ret := libWandMagickReadImageBlob(w.ptr, &blob[0], uint32(len(blob)))
	if !ret {
		var t ExceptionType
		msg := libWandMagickGetException(w.ptr, &t)
		return &MagickException{m: msg, t: t}
	}
	return nil
}

func (w *MagickWand) GetInterlaceScheme() (InterlaceType, error) {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return UndefinedInterlace, ErrInvalidWand
	}
	scheme := libWandMagickGetInterlaceScheme(w.ptr)
	return scheme, nil
}

func (w *MagickWand) SetInterlaceScheme(scheme InterlaceType) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native MagickSetInterlaceScheme
	ret := libWandMagickSetInterlaceScheme(w.ptr, scheme)
	if !ret {
		var t ExceptionType
		msg := libWandMagickGetException(w.ptr, &t)
		return &MagickException{m: msg, t: t}

	}
	return nil
}

func (w *MagickWand) GetImageBlob() ([]byte, error) {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return nil, ErrInvalidWand
	}
	// native GetImageBlob

	var length uint32
	blobPtr := libWandMagickGetImageBlob(w.ptr, &length)
	defer libWandMagickRelinquishMemory(unsafe.Pointer(blobPtr))

	blob := unsafe.Slice(blobPtr, length)[:]

	return blob, nil
}

func (w *MagickWand) SetImageFormat(format string) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native SetImageFormat

	ret := libWandMagickSetImageFormat(w.ptr, format)
	if !ret {
		var t ExceptionType
		msg := libWandMagickGetException(w.ptr, &t)
		return &MagickException{m: msg, t: t}

	}
	return nil
}

// func (w *MagickWand) WriteImageFile(file *os.File) error {
// 	libInit()
// 	if !libWandIsMagickWand(w.ptr) {
// 		return ErrInvalidWand
// 	}
// 	// native MagickWriteImageFile
// 	panic("skeleton: imagick.(MagickWand.WriteImageFile)")
// }

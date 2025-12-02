package imagick

import (
	"errors"
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

// GetLastError returns the kind, reason and description of any error that occurs when using other methods in this API.
// The exception is cleared after this call.
func (mw *MagickWand) GetLastError() error {
	return mw.getLastError(true)
}

// Returns the kind, reason and description of any error that occurs when using other methods in this API.
// Clears the exception, if clear is true.
func (mw *MagickWand) getLastError(clear bool) error {
	var et ExceptionType
	description := libWandMagickGetException(mw.ptr, &et)
	// defer libWandMagickRelinquishMemory(unsafe.Pointer(csdescription))
	if et != UndefinedException {
		if clear {
			mw.clearException()
		}
		return &MagickException{t: et, m: description}
	}
	return nil
}

func (mw *MagickWand) getLastErrorIfFailed(ok bool) error {
	if ok {
		return nil
	}
	return mw.GetLastError()
}

// Clears any exceptions associated with the wand
func (mw *MagickWand) clearException() bool {
	return libWandMagickClearException(mw.ptr)
}

func (w *MagickWand) ReadImageBlob(blob []byte) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	if len(blob) == 0 {
		return errors.New("zero-length blob not permitted")
	}
	// native MagickReadImageBlob
	ok := libWandMagickReadImageBlob(w.ptr, &blob[0], uint32(len(blob)))
	return w.getLastErrorIfFailed(ok)
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
	ok := libWandMagickSetInterlaceScheme(w.ptr, scheme)
	return w.getLastErrorIfFailed(ok)
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

	if err := w.GetLastError(); err != nil {
		return nil, err
	}

	blob := unsafe.Slice(blobPtr, length)

	return blob, nil
}

func (w *MagickWand) GetImageFormat() string {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ""
	}
	return libWandMagickGetImageFormat(w.ptr)
}

func (w *MagickWand) SetImageFormat(format string) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native SetImageFormat

	ok := libWandMagickSetImageFormat(w.ptr, format)
	return w.getLastErrorIfFailed(ok)
}

// func (w *MagickWand) WriteImageFile(file *os.File) error {
// 	libInit()
// 	if !libWandIsMagickWand(w.ptr) {
// 		return ErrInvalidWand
// 	}
// 	// native MagickWriteImageFile
// 	panic("skeleton: imagick.(MagickWand.WriteImageFile)")
// }

func (w *MagickWand) MagickGetImageInterlaceScheme() (InterlaceType, error) {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return UndefinedInterlace, ErrInvalidWand
	}
	scheme := libWandMagickGetImageInterlaceScheme(w.ptr)
	return scheme, nil
}

func (w *MagickWand) SetImageInterlaceScheme(scheme InterlaceType) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native MagickSetInterlaceScheme
	ok := libWandMagickSetImageInterlaceScheme(w.ptr, scheme)
	return w.getLastErrorIfFailed(ok)
}

func (w *MagickWand) WriteImage(filename string) error {
	libInit()
	if !libWandIsMagickWand(w.ptr) {
		return ErrInvalidWand
	}
	// native MagickWriteImage
	ok := libWandMagickWriteImage(w.ptr, filename)
	return w.getLastErrorIfFailed(ok)
}

// func (w *MagickWand) WriteImageFile(file *os.File) error {
// 	libInit()
// 	if !libWandIsMagickWand(w.ptr) {
// 		return ErrInvalidWand
// 	}
// 	// native MagickWriteImageFile
// 	ok := libWandMagickWriteImageFile(w.ptr, file.Fd())
// 	return w.getLastErrorIfFailed(ok)
// }

package vlc

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	initLckOnce sync.Mutex
	initPtr     uintptr
	initError   error
)

func mustGetSymbol(sym string) uintptr {
	ptr, err := getSymbol(sym)
	if err != nil {
		panic(err)
	}
	return ptr
}

func libInitFuncs() {
	purego.RegisterLibFunc(&libvlcGetVersion, initPtr, "libvlc_get_version")
	purego.RegisterLibFunc(&libvlcGetCompiler, initPtr, "libvlc_get_compiler")
	purego.RegisterLibFunc(&libvlcGetChangeset, initPtr, "libvlc_get_changeset")

	purego.RegisterLibFunc(&libvlcErrmsg, initPtr, "libvlc_errmsg")

	purego.RegisterLibFunc(&libvlcNew, initPtr, "libvlc_new")
	purego.RegisterLibFunc(&libvlcRelease, initPtr, "libvlc_release")
	purego.RegisterLibFunc(&libvlcRetain, initPtr, "libvlc_retain")
	purego.RegisterLibFunc(&libvlcFree, initPtr, "libvlc_free")

	purego.RegisterLibFunc(&libvlcMediaNewLocation, initPtr, "libvlc_media_new_location")
	purego.RegisterLibFunc(&libvlcMediaNewPath, initPtr, "libvlc_media_new_path")
	purego.RegisterLibFunc(&libvlcMediaRelease, initPtr, "libvlc_media_release")

	purego.RegisterLibFunc(&libvlcMediaPlayerNewFromMedia, initPtr, "libvlc_media_player_new_from_media")
	purego.RegisterLibFunc(&libvlcMediaPlayerRelease, initPtr, "libvlc_media_player_release")
	purego.RegisterLibFunc(&libvlcMediaPlayerPlay, initPtr, "libvlc_media_player_play")
	purego.RegisterLibFunc(&libvlcMediaPlayerStop, initPtr, "libvlc_media_player_stop")

}

var (
	libvlcGetVersion   func() string
	libvlcGetCompiler  func() string
	libvlcGetChangeset func() string

	libvlcErrmsg func() string

	libvlcNew     func(argc int32, argv unsafe.Pointer) libvlcInstance
	libvlcRelease func(instance libvlcInstance)
	libvlcRetain  func(instance libvlcInstance) libvlcInstance
	libvlcFree    func(ptr unsafe.Pointer)

	libvlcMediaNewLocation func(instance libvlcInstance, location string) libvlcMedia
	libvlcMediaNewPath     func(instance libvlcInstance, path string) libvlcMedia
	libvlcMediaRelease     func(media libvlcMedia)

	libvlcMediaPlayerNewFromMedia func(media libvlcMedia) libvlcMediaPlayer
	libvlcMediaPlayerRelease      func(mediaPlayer libvlcMediaPlayer)
	libvlcMediaPlayerPlay         func(mediaPlayer libvlcMediaPlayer)
	libvlcMediaPlayerStop         func(mediaPlayer libvlcMediaPlayer)
)

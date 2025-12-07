package vlc

import "unsafe"

type libvlcMediaPlayer unsafe.Pointer

type MediaPlayer struct {
	ptr libvlcMediaPlayer
}

func (m *Media) NewMediaPlayer() *MediaPlayer {
	ptr := libvlcMediaPlayerNewFromMedia(m.ptr)
	mp := &MediaPlayer{
		ptr: ptr,
	}
	return mp
}

func (mp *MediaPlayer) Close() error {
	libvlcMediaPlayerRelease(mp.ptr)
	return nil
}

func (mp *MediaPlayer) Play() {
	libvlcMediaPlayerPlay(mp.ptr)
}

func (mp *MediaPlayer) Stop() {
	libvlcMediaPlayerStop(mp.ptr)
}

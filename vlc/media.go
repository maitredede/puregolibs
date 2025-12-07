package vlc

import "unsafe"

type libvlcMedia unsafe.Pointer

type Media struct {
	inst *Instance

	ptr libvlcMedia
}

func (i *Instance) NewMediaFromPath(path string) *Media {
	mPtr := libvlcMediaNewPath(i.ptr, path)
	m := &Media{
		inst: i,
		ptr:  mPtr,
	}
	return m
}

func (i *Instance) NewMediaFromLocation(location string) *Media {
	mPtr := libvlcMediaNewLocation(i.ptr, location)
	m := &Media{
		inst: i,
		ptr:  mPtr,
	}
	return m
}

func (m *Media) Close() error {
	libvlcMediaRelease(m.ptr)
	return nil
}

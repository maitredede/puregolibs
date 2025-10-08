package sane

import "github.com/maitredede/puregolibs/sane/internal"

type Info struct {
	Inexact      bool
	ReloadOpts   bool
	ReloadParams bool
}

func infoFromValue(i internal.SANE_Int) Info {
	v := Info{}
	v.Inexact = (i & infoInexact) != 0
	v.ReloadOpts = (i & infoReloadOptions) != 0
	v.ReloadParams = (i & infoReloadParams) != 0
	return v
}

const (
	infoInexact       internal.SANE_Int = 1
	infoReloadOptions internal.SANE_Int = 2
	infoReloadParams  internal.SANE_Int = 4
)

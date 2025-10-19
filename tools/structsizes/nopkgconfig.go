//go:build linux && cgo && (no_pkgconfig || nopkgconfig)

package main

// #cgo LDFLAGS: -lnfc -lfreefare -lsane -lcec -lusb-1.0 -lfontconfig -lfreetype -lMagickWand -lMagickCore -lplutobook
import "C"

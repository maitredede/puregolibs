//go:build cgo && (no_pkgconfig || nopkgconfig)

package main

// #cgo LDFLAGS: -lnfc -lfreefare -lsane
import "C"

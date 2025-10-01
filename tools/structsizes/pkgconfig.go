//go:build cgo && !no_pkgconfig && !nopkgconfig

package main

// #cgo pkg-config: libnfc libfreefare
import "C"

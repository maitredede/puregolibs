//go:build cgo && !no_pkgconfig && !nopkgconfig

package main

// #cgo pkg-config: libnfc libfreefare sane-backends libcec libusb-1.0
import "C"

package main

import "bytes"

func buildEdid() []byte {
	// Header information
	fixed := []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}
	manufacturerID := []byte{0b0_00001_00, 0b001_00001} // AAA
	productCode := []byte{0x00, 0x00}                   // 0000
	serialNum := []byte{0x01, 0x02, 0x03, 0x04}         // 1234
	manufactureWeek := []byte{0x01}                     // first week
	manufactureYear := []byte{0x23}                     // 2025 (2025-1990=35d=23h)
	edidVersion := []byte{0x01, 0x03}
	header := concat(fixed, manufacturerID, productCode, serialNum, manufactureWeek, manufactureYear, edidVersion)

	// Basic display parameters
	videoInput := []byte{0b1_010_0000} // digital input, 8 bits per color, undefined interface
	horizontalSize := []byte{0}
	verticalSize := []byte{0}
	gamma := []byte{0}
	features := []byte{0b0_0_0_00_0_0_0}
	basicParameters := concat(videoInput, horizontalSize, verticalSize, gamma, features)

	// Chromaticity
	rgBits := []byte{0x00}
	bwBits := []byte{0x00}
	redX := []byte{0x00}
	redY := []byte{0x00}
	greenXY := []byte{0x00, 0x00}
	blueXY := []byte{0x00, 0x00}
	whitePointXY := []byte{0x00, 0x00}
	chromacity := concat(rgBits, bwBits, redX, redY, greenXY, blueXY, whitePointXY)

	// Established timing bitmap. Supported bitmap for (formerly) very common timing modes
	commonTimings := []byte{0xFF, 0xFF, 0b1_0000000}

	// Standard timing information
	resolution := []byte{0x01}
	aspect := []byte{0b00_000000} // 16:10 at 60Hz
	otherTimings := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	standardTimings := concat(resolution, aspect, otherTimings)

	// Display timing descriptor
	descriptors := make([]byte, 125-54+1)

	//extensions
	extensionsCount := []byte{0x00}

	edidNoCheck := concat(header, basicParameters, chromacity, commonTimings, standardTimings, descriptors, extensionsCount)
	sum := 0
	for _, b := range edidNoCheck {
		sum = sum + int(b)
	}
	rem := sum % 256
	checkSum := byte(256 - rem)
	edid := append(edidNoCheck, checkSum)

	return edid
}

func concat(blobs ...[]byte) []byte {
	buff := &bytes.Buffer{}
	for _, b := range blobs {
		if _, err := buff.Write(b); err != nil {
			panic(err)
		}
	}
	return buff.Bytes()
}

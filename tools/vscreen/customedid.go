package main

// GenerateEDID génère un EDID valide pour une résolution donnée
// width: largeur en pixels
// height: hauteur en pixels
// refreshRate: taux de rafraîchissement en Hz (typiquement 60)
func GenerateEDID(width, height, refreshRate int) []byte {
	edid := make([]byte, 128)

	// Header (8 bytes) - signature EDID
	edid[0] = 0x00
	edid[1] = 0xFF
	edid[2] = 0xFF
	edid[3] = 0xFF
	edid[4] = 0xFF
	edid[5] = 0xFF
	edid[6] = 0xFF
	edid[7] = 0x00

	// Manufacturer ID (2 bytes) - "ABC" comme exemple
	// Format: 5 bits par lettre, A=1, B=2, C=3
	edid[8] = 0x04 // 00000 00100
	edid[9] = 0x43 // 010 00011 = B C

	// Product code (2 bytes)
	edid[10] = 0x00
	edid[11] = 0x00

	// Serial number (4 bytes)
	edid[12] = 0x01
	edid[13] = 0x00
	edid[14] = 0x00
	edid[15] = 0x00

	// Week of manufacture (1 byte) - semaine 1
	edid[16] = 0x01

	// Year of manufacture (1 byte) - 2024 = 1990 + 34
	edid[17] = 0x22 // 34 en décimal

	// EDID version (2 bytes) - version 1.3
	edid[18] = 0x01
	edid[19] = 0x03

	// Video input definition (1 byte) - digital
	edid[20] = 0x80

	// Max horizontal image size (cm) - 0 = undefined
	edid[21] = 0x00

	// Max vertical image size (cm) - 0 = undefined
	edid[22] = 0x00

	// Display gamma (1 byte) - 2.2 = (2.2-1.0)*100 + 100 = 220 = 0xDC
	edid[23] = 0x78 // 2.2 gamma

	// Feature support (1 byte)
	edid[24] = 0x0A

	// Chromaticity coordinates (10 bytes) - valeurs standard
	for i := 25; i < 35; i++ {
		edid[i] = 0x00
	}

	// Established timings (3 bytes) - tous à 0
	edid[35] = 0x00
	edid[36] = 0x00
	edid[37] = 0x00

	// Standard timing information (16 bytes) - tous inutilisés
	for i := 38; i < 54; i++ {
		edid[i] = 0x01
		i++
		edid[i] = 0x01
	}

	// Detailed Timing Descriptor #1 (18 bytes) - commence à l'offset 54
	setDetailedTiming(edid[54:72], width, height, refreshRate)

	// Detailed Timing Descriptors #2-4 (54 bytes) - Display Range Limits et dummy
	// Descriptor #2: Display Range Limits
	edid[72] = 0x00
	edid[73] = 0x00
	edid[74] = 0x00
	edid[75] = 0xFD // Display Range Limits tag
	edid[76] = 0x00
	edid[77] = byte(refreshRate - 1) // Min vertical rate
	edid[78] = byte(refreshRate + 1) // Max vertical rate
	edid[79] = 0x1E                  // Min horizontal rate (30 kHz)
	edid[80] = 0x5A                  // Max horizontal rate (90 kHz)
	edid[81] = 0x10                  // Max pixel clock / 10 MHz (160 MHz)
	edid[82] = 0x00                  // No extended timing info
	edid[83] = 0x0A
	for i := 84; i < 90; i++ {
		edid[i] = 0x20
	}

	// Descriptor #3: Monitor name
	edid[90] = 0x00
	edid[91] = 0x00
	edid[92] = 0x00
	edid[93] = 0xFC // Monitor name tag
	edid[94] = 0x00
	copy(edid[95:108], []byte("Custom EDID\n"))

	// Descriptor #4: Dummy
	for i := 108; i < 126; i++ {
		edid[i] = 0x00
	}
	edid[108] = 0x00
	edid[109] = 0x00
	edid[110] = 0x00
	edid[111] = 0x10 // Dummy descriptor

	// Extension flag (1 byte)
	edid[126] = 0x00

	// Checksum (1 byte) - doit être calculé pour que la somme de tous les octets % 256 = 0
	edid[127] = calculateChecksum(edid[:127])

	return edid
}

// setDetailedTiming configure un Detailed Timing Descriptor
func setDetailedTiming(descriptor []byte, width, height, refreshRate int) {
	// Calcul des timings basé sur CVT (Coordinated Video Timings)
	// Paramètres simplifiés pour générer des timings valides

	hBlank := width / 5 // H blanking ~ 20% de la largeur
	vBlank := 30        // V blanking fixe

	hSyncOffset := hBlank / 4
	hSyncWidth := hBlank / 10
	vSyncOffset := 3
	vSyncWidth := 6

	totalH := width + hBlank
	totalV := height + vBlank

	// Pixel clock en Hz, converti en unités de 10kHz
	pixelClock := float64(totalH * totalV * refreshRate)
	pixelClock10kHz := int(pixelClock / 10000.0)

	// Limiter le pixel clock
	if pixelClock10kHz > 0xFFFF {
		pixelClock10kHz = 0xFFFF
	}

	// Pixel clock (2 bytes, little endian, en unités de 10kHz)
	descriptor[0] = byte(pixelClock10kHz & 0xFF)
	descriptor[1] = byte((pixelClock10kHz >> 8) & 0xFF)

	// Horizontal addressable pixels (lower 8 bits)
	descriptor[2] = byte(width & 0xFF)

	// Horizontal blanking (lower 8 bits)
	descriptor[3] = byte(hBlank & 0xFF)

	// Horizontal addressable (upper 4 bits) | Horizontal blanking (upper 4 bits)
	descriptor[4] = byte(((width>>8)&0x0F)<<4 | ((hBlank >> 8) & 0x0F))

	// Vertical addressable lines (lower 8 bits)
	descriptor[5] = byte(height & 0xFF)

	// Vertical blanking (lower 8 bits)
	descriptor[6] = byte(vBlank & 0xFF)

	// Vertical addressable (upper 4 bits) | Vertical blanking (upper 4 bits)
	descriptor[7] = byte(((height>>8)&0x0F)<<4 | ((vBlank >> 8) & 0x0F))

	// Horizontal sync offset (lower 8 bits)
	descriptor[8] = byte(hSyncOffset & 0xFF)

	// Horizontal sync pulse width (lower 8 bits)
	descriptor[9] = byte(hSyncWidth & 0xFF)

	// Vertical sync offset (lower 4 bits) | Vertical sync pulse width (lower 4 bits)
	descriptor[10] = byte((vSyncOffset&0x0F)<<4 | (vSyncWidth & 0x0F))

	// Sync bits (upper 2 bits of each timing)
	descriptor[11] = byte(((hSyncOffset>>8)&0x03)<<6 |
		((hSyncWidth>>8)&0x03)<<4 |
		((vSyncOffset>>4)&0x03)<<2 |
		((vSyncWidth >> 4) & 0x03))

	// Horizontal image size (mm) - 0 = undefined
	descriptor[12] = 0x00
	descriptor[13] = 0x00

	// Horizontal & vertical image size upper bits
	descriptor[14] = 0x00

	// Horizontal border (pixels)
	descriptor[15] = 0x00

	// Vertical border (lines)
	descriptor[16] = 0x00

	// Features bitmap
	// bit 7: interlaced (0)
	// bit 6-5: stereo mode (00)
	// bit 4-3: sync type (11 = digital separate)
	// bit 2: vsync polarity (0)
	// bit 1: hsync polarity (0)
	// bit 0: reserved (0)
	descriptor[17] = 0x18 // Digital separate sync
}

// calculateChecksum calcule le checksum EDID
func calculateChecksum(data []byte) byte {
	var sum int
	for _, b := range data {
		sum += int(b)
	}
	return byte((256 - (sum % 256)) % 256)
}

// func main() {
// 	// Exemple d'utilisation
// 	edid := GenerateEDID(1920, 1080, 60)

// 	fmt.Println("EDID généré (128 octets):")
// 	for i := 0; i < len(edid); i += 16 {
// 		for j := 0; j < 16 && i+j < len(edid); j++ {
// 			fmt.Printf("%02X ", edid[i+j])
// 		}
// 		fmt.Println()
// 	}

// 	// Vérification du checksum
// 	var sum int
// 	for _, b := range edid {
// 		sum += int(b)
// 	}
// 	fmt.Printf("\nChecksum valide: %v (somme mod 256 = %d)\n", sum%256 == 0, sum%256)
// }

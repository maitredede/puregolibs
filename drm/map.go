package drm

// ModeMapDumb set up for mmap of a dumb scanout buffer
type ModeMapDumb struct {
	Handle uint32 // Handle for the object being mapped
	Pad    uint32
	Offset uint64 // Fake offset to use for subsequent mmap call. This is a fixed-size type for 32/64 compatibility
}

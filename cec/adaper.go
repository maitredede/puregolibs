package cec

type nativeAdapter struct {
	path [1024]byte
	comm [1024]byte
}

type Adapter struct {
	Path string
	Comm string
}

type nativeAdapterDescriptor struct {
	path              [1024]byte
	comm              [1024]byte
	vendorID          uint16
	productID         uint16
	firmwareVersion   uint16
	physicalAddress   uint16
	firmwareBuildDate uint32
	adapterType       AdapterType
}

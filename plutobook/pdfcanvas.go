package plutobook

type PDFCanvas struct {
	CanvasBase
}

var _ canvasInterface = (*PDFCanvas)(nil)

func CreatePDFCanvas(filename string, size PageSize) (*PDFCanvas, error) {
	libInit()
	// ptr := libPDFCanvasCreate(int32(width), int32(height), PDFFormat)
	// if ptr == 0 {
	// 	return nil, fmt.Errorf("PDF canvas create failed")
	// }
	// c := &PDFCanvas{
	// 	ptr: ptr,
	// }
	// return c, nil
	panic("TODO: CreatePDFCanvas")
}

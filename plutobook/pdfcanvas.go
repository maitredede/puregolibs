package plutobook

type PDFCanvas struct {
	CanvasBase
}

var _ canvasInterface = (*PDFCanvas)(nil)

// func CreatePDFCanvas(filename string, size PageSize) (*PDFCanvas, error) {
// 	libInit()
// 	ptr := libPDFCanvasCreate(filename, size)
// 	if ptr == 0 {
// 		msg := libGetErrorMessage()
// 		return nil, fmt.Errorf("pdf canvas create failed: %v", msg)
// 	}
// 	c := &PDFCanvas{
// 		CanvasBase: CanvasBase{
// 			ptr: ptr,
// 		},
// 	}
// 	return c, nil
// }

func (c *PDFCanvas) SetMetadata(metadata PdfMetadata, value string) error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libPDFCanvasSetMetadata(c.ptr, metadata, value)
	return nil
}

// func (c *PDFCanvas) SetSize(size PageSize) error {
// 	libInit()
// 	if c.ptr == 0 {
// 		return ErrCanvasIsClosed
// 	}
// 	libPDFCanvasSetSize(c.ptr, size)
// 	return nil
// }

func (c *PDFCanvas) ShowPage() error {
	libInit()
	if c.ptr == 0 {
		return ErrCanvasIsClosed
	}
	libPDFCanvasShowPage(c.ptr)
	return nil
}

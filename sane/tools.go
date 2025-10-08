package sane

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

const (
	opaque8  = uint8(0xff)
	opaque16 = uint16(0xffff)
)

type Frame struct {
	Format       Format // frame format
	Width        int    // width in pixels
	Height       int    // height in pixels
	Channels     int    // number of channels
	Depth        int    // bits per sample
	IsLast       bool   // whether this is the last frame
	bytesPerLine int    // bytes per line, including any padding
	data         []byte // raw data
}

// At returns the sample at coordinates (x,y) for channel ch.
// Note that values are not normalized to the uint16 range,
// so you need to interpret them relative to the color depth.
func (f *Frame) At(x, y, ch int) uint16 {
	switch f.Depth {
	case 1:
		i := f.bytesPerLine*y + f.Channels*(x/8) + ch
		s := (f.data[i] >> uint8(x%8)) & 0x01
		if f.Format == FrameGray {
			// For B&W lineart, 0 is white and 1 is black
			return uint16(s ^ 0x1)
		}
		return uint16(s)
	case 8:
		i := f.bytesPerLine*y + f.Channels*x + ch
		return uint16(f.data[i])
	case 16:
		i := f.bytesPerLine*y + 2*(f.Channels*x+ch)
		return uint16(f.data[i+1])<<8 + uint16(f.data[i])
	}
	return 0
}

func (h *Handle) ReadFrame() (*Frame, error) {
	if err := h.Start(); err != nil {
		return nil, err
	}
	p, err := h.GetParameters()
	if err != nil {
		return nil, err
	}
	if p.Depth != 1 && p.Depth != 8 && p.Depth != 16 {
		return nil, fmt.Errorf("unsupported bit depth: %d", p.Depth)
	}

	data := new(bytes.Buffer)
	if p.Lines > 0 {
		// Preallocate buffer with expected size
		data.Grow(p.Lines * p.BytesPerLine)
	}
	n, err := data.ReadFrom(h)
	if err != nil {
		return nil, err
	}
	_ = n
	// buff := make([]byte, 1024)
	// for {
	// 	n, err := Read(h, buff)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		return nil, err
	// 	}
	// 	w, err := data.Write(buff[:n])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if n != w {
	// 		panic("FIXME")
	// 	}
	// }

	nch := 1
	if p.Format == FrameRGB {
		nch = 3
	}
	f := &Frame{
		Format:       p.Format,
		Width:        p.PixelsPerLine,
		Height:       data.Len() / p.BytesPerLine, // p.Lines is unreliable
		Channels:     nch,
		Depth:        p.Depth,
		IsLast:       p.IsLastFrame,
		bytesPerLine: p.BytesPerLine,
		data:         data.Bytes(),
	}
	return f, nil
}

type Image struct {
	fs [3]*Frame // multiple frames must be in RGB order
}

// Bounds returns the domain for which At returns valid pixels.
func (m *Image) Bounds() image.Rectangle {
	f := m.fs[0]
	return image.Rect(0, 0, f.Width, f.Height)
}

// ColorModel returns the Image's color model.
func (m *Image) ColorModel() color.Model {
	f := m.fs[0]
	switch {
	case f.Depth != 16 && f.Format == FrameGray:
		return color.GrayModel
	case f.Depth == 16 && f.Format == FrameGray:
		return color.Gray16Model
	case f.Depth != 16 && f.Format != FrameGray:
		return color.RGBAModel
	case f.Depth == 16 && f.Format != FrameGray:
		return color.RGBA64Model
	}
	return color.RGBAModel
}

// At returns the color of the pixel at (x, y).
func (m *Image) At(x, y int) color.Color {
	if x < 0 || x >= m.fs[0].Width || y < 0 || y >= m.fs[0].Height {
		return color.RGBA{}
	}
	if m.fs[0].Format == FrameGray {
		// grayscale
		switch m.fs[0].Depth {
		case 1:
			return color.Gray{uint8(0xFF * m.fs[0].At(x, y, 0))}
		case 8:
			return color.Gray{uint8(m.fs[0].At(x, y, 0))}
		case 16:
			return color.Gray16{m.fs[0].At(x, y, 0)}
		}
	} else {
		// color
		var r, g, b uint16
		if m.fs[0].Format == FrameRGB {
			// interleaved
			r = m.fs[0].At(x, y, 0)
			g = m.fs[0].At(x, y, 1)
			b = m.fs[0].At(x, y, 2)
		} else {
			// non-interleaved
			r = m.fs[0].At(x, y, 0)
			g = m.fs[1].At(x, y, 0)
			b = m.fs[2].At(x, y, 0)
		}
		switch m.fs[0].Depth {
		case 1:
			return color.RGBA{uint8(0xFF * r), uint8(0xFF * g), uint8(0xFF * b), opaque8}
		case 8:
			return color.RGBA{uint8(r), uint8(g), uint8(b), opaque8}
		case 16:
			return color.RGBA64{r, g, b, opaque16}
		}
	}
	return color.RGBA{} // shouldn't happen
}

func (h *Handle) ReadImage() (*Image, error) {
	defer h.Cancel()

	m := Image{}
	for {
		f, err := h.ReadFrame()
		if err != nil {
			return nil, err
		}
		switch f.Format {
		case FrameGray, FrameRGB, FrameRed:
			m.fs[0] = f
		case FrameGreen:
			m.fs[1] = f
		case FrameBlue:
			m.fs[2] = f
		default:
			return nil, fmt.Errorf("unknown frame type %d", f.Format)
		}
		if f.IsLast {
			break
		}
	}
	return &m, nil
}

// ReadAvailableImages reads all available image from the connection.
// This is required for example for duplex scanners like the Fujitsu
// ix500 as ReadImage only fetches one page from the scanner.
func (h *Handle) ReadAvailableImages() ([]*Image, error) {
	defer h.Cancel()

	var (
		images       = []*Image{}
		keepFetching = true
	)

	for keepFetching {
		m := Image{}
		for {
			f, err := h.ReadFrame()
			if err != nil {
				// if err == ErrEmpty && len(images) > 0 {
				if se, ok := err.(*saneStatusError); ok && se.s == StatusNoDocs && len(images) > 0 {
					// This is expected in multi-page scenarios and signals
					// there are no more pages to come.
					keepFetching = false
					break
				}

				// Other errors are returned
				return nil, err
			}
			switch f.Format {
			case FrameGray, FrameRGB, FrameRed:
				m.fs[0] = f
			case FrameGreen:
				m.fs[1] = f
			case FrameBlue:
				m.fs[2] = f
			default:
				return nil, fmt.Errorf("unknown frame type %d", f.Format)
			}
			if f.IsLast {
				images = append(images, &m)
				keepFetching = false
				break
			}
		}
	}

	return images, nil
}

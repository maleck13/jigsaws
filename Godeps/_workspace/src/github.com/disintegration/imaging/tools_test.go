package imaging

import (
	"image"
	"testing"
)

func TestCrop(t *testing.T) {
	td := []struct {
		desc string
		src  image.Image
		r    image.Rectangle
		want *image.NRGBA
	}{
		{
			"Crop 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			image.Rect(-1, 0, 1, 1),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
				},
			},
		},
	}
	for _, d := range td {
		got := Crop(d.src, d.r)
		want := d.want
		if !compareNRGBA(got, want, 0) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

func TestCropCenter(t *testing.T) {
	td := []struct {
		desc string
		src  image.Image
		w, h int
		want *image.NRGBA
	}{
		{
			"CropCenter 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			2, 1,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
				},
			},
		},
	}
	for _, d := range td {
		got := CropCenter(d.src, d.w, d.h)
		want := d.want
		if !compareNRGBA(got, want, 0) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

func TestPaste(t *testing.T) {
	td := []struct {
		desc string
		src1 image.Image
		src2 image.Image
		p    image.Point
		want *image.NRGBA
	}{
		{
			"Paste 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				},
			},
			image.Pt(-1, 0),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
	}
	for _, d := range td {
		got := Paste(d.src1, d.src2, d.p)
		want := d.want
		if !compareNRGBA(got, want, 0) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

func TestPasteCenter(t *testing.T) {
	td := []struct {
		desc string
		src1 image.Image
		src2 image.Image
		want *image.NRGBA
	}{
		{
			"PasteCenter 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
	}
	for _, d := range td {
		got := PasteCenter(d.src1, d.src2)
		want := d.want
		if !compareNRGBA(got, want, 0) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

func TestOverlay(t *testing.T) {
	td := []struct {
		desc string
		src1 image.Image
		src2 image.Image
		p    image.Point
		a    float64
		want *image.NRGBA
	}{
		{
			"Overlay 2x3 2x1 1.0",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x60, 0x00, 0x90, 0xff, 0xff, 0x00, 0x99, 0x7f,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x20, 0x40, 0x80, 0x7f, 0xaa, 0xbb, 0xcc, 0xff,
				},
			},
			image.Pt(-1, 0),
			1.0,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x40, 0x1f, 0x88, 0xff, 0xaa, 0xbb, 0xcc, 0xff,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
		{
			"Overlay 2x2 2x2 0.5",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x00, 0x00, 0xff, 0xff, 0x20, 0x20, 0x20, 0x00,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
					0xff, 0xff, 0x00, 0xff, 0x20, 0x20, 0x20, 0xff,
				},
			},
			image.Pt(-1, -1),
			0.5,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x7f, 0x7f, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x7f, 0x7f, 0x7f, 0xff, 0x20, 0x20, 0x20, 0x7f,
				},
			},
		},
	}
	for _, d := range td {
		got := Overlay(d.src1, d.src2, d.p, d.a)
		want := d.want
		if !compareNRGBA(got, want, 1) {
			t.Errorf("test [%s] failed: %#v", d.desc, got)
		}
	}
}

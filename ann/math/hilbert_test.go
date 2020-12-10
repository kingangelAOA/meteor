package math

import "testing"

import "github.com/fogleman/gg"

func TestHilbert(t *testing.T)  {
	hilbert(0, 0, width, 0, 0)
	dc := gg.NewContext(650, 650)
	dc.SetRGB(0, 0, 0) // Black background
	dc.Clear()
	for _, p := range points {
		dc.LineTo(p.X, p.Y)
	}
	dc.SetHexColor("#90EE90") // Light green curve
	dc.SetLineWidth(1)
	dc.Stroke()
	dc.SavePNG("hilbert.png")
}

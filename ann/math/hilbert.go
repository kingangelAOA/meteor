package math

import "github.com/fogleman/gg"

var points []gg.Point

const width = 4

func hilbert(x, y, lg, i1, i2 int) {
	if lg == 1 {
		px := float64(width-x) * 10
		py := float64(width-y) * 10
		points = append(points, gg.Point{px, py})
		return
	}
	lg >>= 1
	hilbert(x+i1*lg, y+i1*lg, lg, i1, 1-i2)
	hilbert(x+i2*lg, y+(1-i2)*lg, lg, i1, i2)
	hilbert(x+(1-i1)*lg, y+(1-i1)*lg, lg, i1, i2)
	hilbert(x+(1-i2)*lg, y+i2*lg, lg, 1-i1, i2)
}

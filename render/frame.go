package render

const (
	WIDTH  = 256
	HEIGHT = 240
)

type Frame struct {
	Screen [WIDTH * HEIGHT * 3]uint8
	//essentially a flattened rep of a 3 dimensional image x,y and 3 values for rgb
}

func (f *Frame) SetPixel(x int, y int, rgb []uint8) {
	base := y*WIDTH*3 + x*3
	if base+2 < len(f.Screen) {
		f.Screen[base] = rgb[0]
		f.Screen[base+1] = rgb[1]
		f.Screen[base+2] = rgb[2]
	}
}

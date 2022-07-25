package render

const (
	WIDTH  = 256
	HEIGHT = 240
)

type Frame struct {
	Screen [256 * 240 * 3]uint8
}

func (f *Frame) SetPixel(x int, y int, rgb [3]uint8) {
	base := y*WIDTH*3 + x*3
	if base+2 < len(f.Screen) {
		f.Screen[base] = rgb[0]
		f.Screen[base+1] = rgb[1]
		f.Screen[base+2] = rgb[2]
	}
}

func (f *Frame) ShowTiles(rom []byte) {

	col0 := [3]uint8{0x00, 0x3d, 0xa6}
	col1 := [3]uint8{0xd6, 0x32, 0x00}
	col2 := [3]uint8{0x00, 0x8a, 0x55}
	col3 := [3]uint8{0x09, 0x09, 0x09}

	//4kb banks hence the multiplication

	//each tile is an 8x8 box with each line repesenting 8 bits aka 1 byte hence 1 tile is 8bytes.Color
	//info is tored in the next occuring 8 bytes for the earloer 8 bytes giving a total of 16 bytes
	//if my hypothesis of 4kb pages is right we only have 2 pages for 8kb of data aka 512 tiles at at 128 bits (colorn included)
	width := 0
	height := 0
	//SCREEN OFFSETS
	bank := 0
	tileNum := 0
	for true {

		tile := rom[((bank * 0x1000) + tileNum*16) : (bank*0x1000)+tileNum*16+16]
		for y := 0; y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				col := upper&1<<1 | lower&1
				upper = upper >> 1
				lower = lower >> 1
				switch col {
				case 0:
					f.SetPixel((int(width*8))+x, (int(height*8))+y, col0)
				case 1:
					f.SetPixel((int(width*8))+x, (int(height*8))+y, col1)
				case 2:
					f.SetPixel((int(width*8))+x, (int(height*8))+y, col2)
				case 3:
					f.SetPixel((int(width*8))+x, (int(height*8))+y, col3)
				}
			}
		}
		tileNum++
		if tileNum == 256 {
			tileNum = 0
			bank++
		}
		if bank == 2 {
			break
		}
		width++
		if width == 32 {
			height++
			width = 0
		}

	}
}

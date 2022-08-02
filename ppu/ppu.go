package ppu

import (
	"emulator/common"
	"emulator/render"
)

const (
	BASE_NAME_TABLE_ONE = iota
	BASE_NAME_TABLE_TWO
	VRAM_INCREMENT
	SPRITE_ADDRESS
	BACKGROUND_PATTERN
	SPRITE_SIZE
	PPU_MASTER
	GENERATE_NMI
	TOP_LEFT
	TOP_RIGHT
	BOTTOM_LEFT
	BOTTOM_RIGHT
	CHR_ROM_START = 0
	CHR_ROM_END   = 0x1FFF
	PPU_RAM_START = 0x2000
	PPU_RAM_END   = 0x3EFF
	PALETTE_START = 0x3F00
	PALETTE_END   = 0x3FFF
)

var pallete = [][]uint8{
	{0x80, 0x80, 0x80}, {0x00, 0x3D, 0xA6}, {0x00, 0x12, 0xB0}, {0x44, 0x00, 0x96}, {0xA1, 0x00, 0x5E},
	{0xC7, 0x00, 0x28}, {0xBA, 0x06, 0x00}, {0x8C, 0x17, 0x00}, {0x5C, 0x2F, 0x00}, {0x10, 0x45, 0x00},
	{0x05, 0x4A, 0x00}, {0x00, 0x47, 0x2E}, {0x00, 0x41, 0x66}, {0x00, 0x00, 0x00}, {0x05, 0x05, 0x05},
	{0x05, 0x05, 0x05}, {0xC7, 0xC7, 0xC7}, {0x00, 0x77, 0xFF}, {0x21, 0x55, 0xFF}, {0x82, 0x37, 0xFA},
	{0xEB, 0x2F, 0xB5}, {0xFF, 0x29, 0x50}, {0xFF, 0x22, 0x00}, {0xD6, 0x32, 0x00}, {0xC4, 0x62, 0x00},
	{0x35, 0x80, 0x00}, {0x05, 0x8F, 0x00}, {0x00, 0x8A, 0x55}, {0x00, 0x99, 0xCC}, {0x21, 0x21, 0x21},
	{0x09, 0x09, 0x09}, {0x09, 0x09, 0x09}, {0xFF, 0xFF, 0xFF}, {0x0F, 0xD7, 0xFF}, {0x69, 0xA2, 0xFF},
	{0xD4, 0x80, 0xFF}, {0xFF, 0x45, 0xF3}, {0xFF, 0x61, 0x8B}, {0xFF, 0x88, 0x33}, {0xFF, 0x9C, 0x12},
	{0xFA, 0xBC, 0x20}, {0x9F, 0xE3, 0x0E}, {0x2B, 0xF0, 0x35}, {0x0C, 0xF0, 0xA4}, {0x05, 0xFB, 0xFF},
	{0x5E, 0x5E, 0x5E}, {0x0D, 0x0D, 0x0D}, {0x0D, 0x0D, 0x0D}, {0xFF, 0xFF, 0xFF}, {0xA6, 0xFC, 0xFF},
	{0xB3, 0xEC, 0xFF}, {0xDA, 0xAB, 0xEB}, {0xFF, 0xA8, 0xF9}, {0xFF, 0xAB, 0xB3}, {0xFF, 0xD2, 0xB0},
	{0xFF, 0xEF, 0xA6}, {0xFF, 0xF7, 0x9C}, {0xD7, 0xE8, 0x95}, {0xA6, 0xED, 0xAF}, {0xA2, 0xF2, 0xDA},
	{0x99, 0xFF, 0xFC}, {0xDD, 0xDD, 0xDD}, {0x11, 0x11, 0x11}, {0x11, 0x11, 0x11},
}

type Ppu struct {
	ChrRom []uint8
	//from rom
	Palette [32]uint8
	//colors
	Ram [2048]uint8
	//ppu mem
	Oam [256]uint8
	//sprite state monitoring
	Mirror int

	AddrRegister addrReg
	OamAddr      PPU_OAM_ADDRESS
	OamData      PPU_OAM_DATA
	Status       PPU_STATUS_REGISTER
	Scroll       PPU_SCROLL_REGISTER
	Mask         PPU_MASK
	//DataRegister    dataReg
	ControlRegister PPU_CONTROL
	NmiOcurred      bool

	buffer    uint8
	PpuTicks  int
	Scanlines int
	Frame     render.Frame
}

func screenToAttribute(width int, height int) int {
	//4x4 tiles of screen need to be chunked together
	//we have 32 tiles in the wifth axis making it chunkable 8 times
	height = (height / 4) * 8
	return height + width/4
}

func (ppu *Ppu) ShowTiles() {

	//4kb banks hence the multiplication

	//each tile is an 8x8 box with each line repesenting 8 bits aka 1 byte hence 1 tile is 8bytes.Color
	//info is tored in the next occuring 8 bytes for the earloer 8 bytes giving a total of 16 bytes
	//if my hypothesis of 4kb pages is right we only have 2 pages for 8kb of data aka 512 tiles at at 128 bits (color included)
	width := 0
	height := 0
	//SCREEN OFFSETS

	for i := 0; i <= 0x3c0; i++ {
		bank := ppu.ControlRegister.GetBackgroundTableAddress()
		tileNum := uint16(ppu.Ram[i])

		tile := ppu.ChrRom[((bank) + tileNum*16) : (bank)+tileNum*16+16]
		attributeIndex := screenToAttribute(width, height)
		colorIndex := ppu.Ram[0x3c0+uint16(attributeIndex)]
		rowOrientation := ((colorIndex * 4) % 32) + 1
		colOrientation := ((colorIndex * 4 * 4) / 32) + 1

		color := uint8(0)

		switch {
		case width <= int(rowOrientation) && height <= int(colOrientation):
			color = colorIndex & 0b00000011
		case width > int(rowOrientation) && height <= int(colOrientation):
			color = (colorIndex >> 2) & 0b00000011
		case width <= int(rowOrientation) && height > int(colOrientation):
			color = (colorIndex >> 4) & 0b00000011
		case width > int(rowOrientation) && height > int(colOrientation):
			color = (colorIndex >> 6) & 0b00000011
		default:
			panic("error!")
		}

		color = 1 + (color * 4)
		colors := []uint8{ppu.Palette[0], ppu.Palette[color], ppu.Palette[color+1], ppu.Palette[color+2]}
		for y := 0; y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				col := upper&1<<1 | lower&1
				upper = upper >> 1
				lower = lower >> 1
				screenRow := (int(width * 8)) + x
				screenCol := (int(height * 8)) + y
				switch col {
				case 0:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[0]])
				case 1:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[1]])
				case 2:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[2]])
				case 3:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[3]])
				}
			}
		}
		width++
		if width == 32 {
			height++
			width = 0
		}

	}
	for i := 0; i < len(ppu.Oam); i = i + 4 {
		yLoc := ppu.Oam[i]
		tileNum := uint16(ppu.Oam[i+1])
		bank := ppu.ControlRegister.GetSpriteTableAddress()
		tile := ppu.ChrRom[((bank) + tileNum*16) : (bank)+tileNum*16+16]
		xLoc := ppu.Oam[i+3]
		attr := ppu.Oam[i+2]

		verticalFlip := false
		horizontalFlip := false

		palleteNum := attr & 0b11

		if attr>>7&0b1 != 0 {
			verticalFlip = true
		}
		if attr>>6&0b01 != 0 {
			horizontalFlip = true
		}
		color := 1 + (palleteNum * 4)
		//17 to skip background colors?????? wrong
		colors := []uint8{0, ppu.Palette[color], ppu.Palette[color+1], ppu.Palette[color+2]}

		for y := 0; y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				col := upper&1<<1 | lower&1
				upper = upper >> 1
				lower = lower >> 1
				screenRow := 0
				screenCol := 0
				switch {
				case horizontalFlip && verticalFlip:
					screenRow = 7 + int(xLoc) - x
					screenCol = 7 - y + int(yLoc)
				case !horizontalFlip && !verticalFlip:
					screenRow = x + int(xLoc)
					screenCol = y + int(yLoc)
				case !horizontalFlip && verticalFlip:
					screenRow = x + int(xLoc)
					screenCol = 7 - y + int(yLoc)
				case horizontalFlip && !verticalFlip:
					screenRow = 7 + int(xLoc) - x
					screenCol = y + int(yLoc)
				default:
					panic("impossible")
				}
				switch col {
				case 1:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[1]])
				case 2:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[2]])
				case 3:
					ppu.Frame.SetPixel(screenRow, screenCol, pallete[colors[3]])
				}
			}
		}

	}
}

func (ppu *Ppu) Tick(amount int) bool {
	ppu.PpuTicks += amount
	if ppu.PpuTicks >= 341 {

		ppu.PpuTicks -= 341
		ppu.Scanlines++
		//fmt.Println(ppu.Scanlines)
		if ppu.Scanlines == 241 {

			ppu.Status.SetVBlank()
			ppu.Status.ClearSpriteZero()
			if ppu.ControlRegister.GenerateNmi() {
				ppu.NmiOcurred = true
			}
		}
		if ppu.Scanlines >= 262 {
			ppu.Status.ClearSpriteZero()
			ppu.Scanlines = 0
			ppu.Status.ClearVBlank()
			ppu.NmiOcurred = false
			return true
		}
	}
	return false
}

//do you want circular refrences? lol
func (ppu *Ppu) PollNmi() bool {
	return ppu.NmiOcurred
}

func NewPPU(rom []uint8, mirror int) *Ppu {

	ppu := &Ppu{
		Mirror: mirror,
		ChrRom: rom,
	}
	//	ppu.PpuTicks = 7 * 3

	return ppu
}

type PPU_OAM_ADDRESS uint8

func (reg *PPU_OAM_ADDRESS) WriteAddressOam(addr uint8) {
	*reg = PPU_OAM_ADDRESS(addr)
}

func (reg *PPU_OAM_ADDRESS) Increment() {
	*reg++
}

func (ppu *Ppu) WriteDataOam(data uint8) {
	ppu.Oam[ppu.OamAddr] = data
	ppu.OamAddr.Increment()
}
func (ppu *Ppu) ReadDataOamRegister() uint8 {
	return ppu.Oam[ppu.OamAddr]
}

//sus
func (ppu *Ppu) WriteOamDMA(data []uint8) {
	//fmt.Println(data)
	for i := 0; i < len(data); i++ {
		ppu.Oam[ppu.OamAddr] = data[i]
		ppu.OamAddr.Increment()
	}

}

type PPU_OAM_DATA uint8

type PPU_MASK uint8

const (
	RED = iota
	BLUE
	GREEN
)

func (ctrl *PPU_MASK) Update(val uint8) {
	*ctrl = PPU_MASK(val)
}

func (ctrl *PPU_MASK) BackgroundRender() bool {
	return hasBit(uint8(*ctrl), 3)
}
func (ctrl *PPU_MASK) SpriteRender() bool {
	return hasBit(uint8(*ctrl), 4)
}

func (ctrl *PPU_MASK) BackgroundRenderTop() bool {
	return hasBit(uint8(*ctrl), 1)
}
func (ctrl *PPU_MASK) SpriteRenderTop() bool {
	return hasBit(uint8(*ctrl), 2)
}

func (ctrl *PPU_MASK) IsGreyScale() bool {
	return hasBit(uint8(*ctrl), 0)
}

func (ctrl *PPU_MASK) EmphasizeRed() bool {
	return hasBit(uint8(*ctrl), 5)
}
func (ctrl *PPU_MASK) EmphasizeGreen() bool {
	return hasBit(uint8(*ctrl), 6)
}

func (ctrl *PPU_MASK) EmphasizeBlue() bool {
	return hasBit(uint8(*ctrl), 7)
}

func (ctrl *PPU_MASK) EnableRendring() bool {
	if uint8(*ctrl) == 0x1e {
		return true
	}

	if uint8(*ctrl) == 0x00 {
		return false
	}
	return true
}

type addrReg struct {
	values [2]uint8
	ptr    int
}

func (reg *addrReg) Update(value uint8) {
	reg.values[reg.ptr] = value
	reg.ptr++
	reg.ptr = (reg.ptr) % 2
	if reg.Get() > 0x3fff {
		reg.Set(reg.Get() & 0x3fff)
		//mirror back to ppu registers
	}
}

func (reg *addrReg) Get() uint16 {
	return (uint16(reg.values[0]))<<8 | (uint16(reg.values[1]))
}

func (reg *addrReg) Set(val uint16) {
	hi := uint8(val >> 8)
	low := uint8(val & 0x00FF)
	reg.values[0] = hi
	reg.values[1] = low
}

func (reg *addrReg) Increment(val uint8) {
	reg.Set(reg.Get() + uint16(val))
}

type PPU_CONTROL uint8

func (ctrl *PPU_CONTROL) ValueToIncrementBy() uint8 {
	if hasBit(uint8(*ctrl), VRAM_INCREMENT) {
		return 32
	} else {
		return 1
	}
}

func (ctrl *PPU_CONTROL) Update(val uint8) {
	*ctrl = PPU_CONTROL(val)
}
func (ctrl *PPU_CONTROL) GetBaseNameTableAddress() uint16 {
	a := getBit(uint8(*ctrl), 0)
	b := getBit(uint8(*ctrl), 1)

	switch {
	case a == 0 && b == 0:
		return 0x2000
	case a == 0 && b == 1:
		return 0x2400
	case a == 1 && b == 0:
		return 0x2800
	case a == 1 && b == 1:
		return 0x2c00
	default:
		panic("error!")
	}
}
func (ctrl *PPU_CONTROL) GetSpriteTableAddress() uint16 {
	a := getBit(uint8(*ctrl), 3)

	switch {
	case a == 0:
		return 0
	case a == 1:
		return 0x1000
	default:
		panic("impossible")
	}

}
func (ctrl *PPU_CONTROL) GetBackgroundTableAddress() uint16 {
	a := getBit(uint8(*ctrl), 4)

	switch {
	case a == 0:
		return 0
	case a == 1:
		return 0x1000
	default:
		panic("impossible")
	}

}
func (ctrl *PPU_CONTROL) GetSpritesize() uint8 {
	a := getBit(uint8(*ctrl), 5)

	switch {
	case a == 0:
		return 8
	case a == 1:
		return 16
	default:
		panic("impossible")
	}

}

func (ctrl *PPU_CONTROL) GetMasterSlave() uint8 {
	return getBit(uint8(*ctrl), 6)
}
func (ctrl *PPU_CONTROL) GenerateNmi() bool {
	return hasBit(uint8(*ctrl), 7)
}

type PPU_SCROLL_REGISTER struct {
	values [2]uint8
	ptr    int
}

func (reg *PPU_SCROLL_REGISTER) Update(value uint8) {
	reg.values[reg.ptr] = value
	reg.ptr++
	reg.ptr = (reg.ptr) % 2
}

type PPU_STATUS_REGISTER uint8

func (reg *PPU_STATUS_REGISTER) SetVBlank() {
	*reg = PPU_STATUS_REGISTER(setBit(uint8(*reg), 7))
}
func (reg *PPU_STATUS_REGISTER) ClearVBlank() {
	*reg = PPU_STATUS_REGISTER(clearBit(uint8(*reg), 7))
}

func (reg *PPU_STATUS_REGISTER) SetSpriteZero() {
	*reg = PPU_STATUS_REGISTER(setBit(uint8(*reg), 6))
}
func (reg *PPU_STATUS_REGISTER) ClearSpriteZero() {
	*reg = PPU_STATUS_REGISTER(clearBit(uint8(*reg), 6))
}

func (reg *PPU_STATUS_REGISTER) SetSpriteOverflow() {
	*reg = PPU_STATUS_REGISTER(setBit(uint8(*reg), 5))
}
func (reg *PPU_STATUS_REGISTER) ClearSpriteOverflow() {
	*reg = PPU_STATUS_REGISTER(clearBit(uint8(*reg), 5))
}
func (reg *PPU_STATUS_REGISTER) InVBlank() bool {
	return hasBit(uint8(*reg), 7)
}

//sus
func (ppu *Ppu) mirriorPPU(addr uint16) uint16 {
	if addr >= 0x3000 && addr <= 0x3eff {
		addr = addr & 0x2fff
	}
	mirror := ppu.Mirror
	switch {
	case mirror == common.HORIZONTAL:
		if addr >= 0x2000 && addr < 0x2400 {
			return addr - 0x2000
		}
		if addr >= 0x2400 && addr < 0x2800 {
			return addr - 0x2400
		}

		if addr >= 0x2800 && addr < 0x2c00 {
			return (addr - 0x2800) + 0x400
		}
		if addr >= 0x2c00 && addr < 0x3f00 {
			return (addr - 0x2c00) + 0x400
		}
	case mirror == common.VERTICAL:
		if addr >= 0x2000 && addr < 0x2400 {
			return addr - 0x2000
		}
		if addr >= 0x2400 && addr < 0x2800 {
			return addr - 0x2400 + 0x400
		}

		if addr >= 0x2800 && addr < 0x2c00 {
			return (addr - 0x2800)
		}
		if addr >= 0x2c00 && addr < 0x3f00 {
			return (addr - 0x2c00) + 0x400
		}
	}
	return 0
}

func (ppu *Ppu) ReadData() uint8 {

	addr := ppu.AddrRegister.Get()
	switch addr {
	case 0x3f10, 0x3f14, 0x3f18, 0x3f1c:
		addr = addr - 0x10
	}
	val := ppu.ControlRegister.ValueToIncrementBy()
	ppu.AddrRegister.Increment(val)

	switch {

	case addr >= PALETTE_START && addr <= PALETTE_END:
		//fmt.Println(ppu.Palette)
		return ppu.Palette[(addr - PALETTE_START)]
	case addr >= CHR_ROM_START && addr <= CHR_ROM_END:
		result := ppu.buffer
		ppu.buffer = ppu.ChrRom[(addr)]
		return result
	case addr >= PPU_RAM_START && addr <= PPU_RAM_END:
		result := ppu.buffer
		ppu.buffer = ppu.Ram[(ppu.mirriorPPU(addr))]
		return result
	}
	return 0
}

func (ppu *Ppu) WriteData(data uint8) {
	addr := ppu.AddrRegister.Get()
	switch addr {
	case 0x3f10, 0x3f14, 0x3f18, 0x3f1c:
		addr = addr - 0x10
	}

	switch {
	case addr >= PALETTE_START && addr <= PALETTE_END:
		ppu.Palette[(addr - PALETTE_START)] = data
	case addr >= PPU_RAM_START && addr <= PPU_RAM_END:
		ppu.Ram[(ppu.mirriorPPU(addr))] = data
	case addr >= CHR_ROM_START && addr <= CHR_ROM_END:
		ppu.ChrRom[(addr)] = data
	}

	val := ppu.ControlRegister.ValueToIncrementBy()
	ppu.AddrRegister.Increment(val)

}
func (ppu *Ppu) WriteToCtrl(val uint8) {
	bitsBefore := uint8(ppu.ControlRegister)
	ppu.ControlRegister.Update(val)
	bitsAfter := uint8(ppu.ControlRegister)
	if !hasBit(bitsBefore, 7) && hasBit(bitsAfter, 7) && ppu.Status.InVBlank() {
		ppu.NmiOcurred = true
	}
	//if we were in vblank but control says we cannot generate an nmi, if we decide to set control to generate nmi while maintaining vblank status
	//we should notify we are already in nmi
}
func (ppu *Ppu) ReadStatus() uint8 {
	temp := uint8(ppu.Status)
	ppu.Status.ClearVBlank()
	ppu.Scroll.ptr = 0
	ppu.AddrRegister.ptr = 0
	return temp
}

//this is a duplicate remove later
func hasBit(n uint8, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func setBit(num uint8, pos int) uint8 {

	num |= (uint8(1) << pos)
	return num
}

func clearBit(n uint8, pos int) uint8 {
	var mask uint8 = ^(1 << pos)
	n &= mask
	return n
}

func getBit(n uint8, pos int) uint8 {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}

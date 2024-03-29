package ppu

import (
	"github.com/Saf1u/bareNES/common"
	"github.com/Saf1u/bareNES/render"
	"github.com/Saf1u/bareNES/utils"
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

func NewPPU(rom []uint8, mirror int) *Ppu {

	ppu := &Ppu{
		Mirror: mirror,
		ChrRom: rom,
	}

	return ppu
}

type PPU_OAM_ADDRESS uint8

type PPU_OAM_DATA uint8

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

func (ppu *Ppu) WriteOamDMA(data []uint8) {
	for i := 0; i < len(data); i++ {
		ppu.Oam[ppu.OamAddr] = data[i]
		ppu.OamAddr.Increment()
	}

}

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
	return utils.HasBit(uint8(*ctrl), 3)
}
func (ctrl *PPU_MASK) SpriteRender() bool {
	return utils.HasBit(uint8(*ctrl), 4)
}

func (ctrl *PPU_MASK) BackgroundRenderTop() bool {
	return utils.HasBit(uint8(*ctrl), 1)
}
func (ctrl *PPU_MASK) SpriteRenderTop() bool {
	return utils.HasBit(uint8(*ctrl), 2)
}

func (ctrl *PPU_MASK) IsGreyScale() bool {
	return utils.HasBit(uint8(*ctrl), 0)
}

func (ctrl *PPU_MASK) EmphasizeRed() bool {
	return utils.HasBit(uint8(*ctrl), 5)
}
func (ctrl *PPU_MASK) EmphasizeGreen() bool {
	return utils.HasBit(uint8(*ctrl), 6)
}

func (ctrl *PPU_MASK) EmphasizeBlue() bool {
	return utils.HasBit(uint8(*ctrl), 7)
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
	if utils.HasBit(uint8(*ctrl), VRAM_INCREMENT) {
		return 32
	} else {
		return 1
	}
}

func (ctrl *PPU_CONTROL) Update(val uint8) {
	*ctrl = PPU_CONTROL(val)
}
func (ctrl *PPU_CONTROL) GetBaseNameTableAddress() uint16 {
	a := utils.GetBit(uint8(*ctrl), 0)
	b := utils.GetBit(uint8(*ctrl), 1)

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
	a := utils.GetBit(uint8(*ctrl), 3)

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
	a := utils.GetBit(uint8(*ctrl), 4)

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
	a := utils.GetBit(uint8(*ctrl), 5)

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
	return utils.GetBit(uint8(*ctrl), 6)
}
func (ctrl *PPU_CONTROL) GenerateNmi() bool {
	return utils.HasBit(uint8(*ctrl), 7)
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
	*reg = PPU_STATUS_REGISTER(utils.SetBit(uint8(*reg), 7))
}
func (reg *PPU_STATUS_REGISTER) ClearVBlank() {
	*reg = PPU_STATUS_REGISTER(utils.ClearBit(uint8(*reg), 7))
}

func (reg *PPU_STATUS_REGISTER) SetSpriteZero() {
	*reg = PPU_STATUS_REGISTER(utils.SetBit(uint8(*reg), 6))
}
func (reg *PPU_STATUS_REGISTER) ClearSpriteZero() {
	*reg = PPU_STATUS_REGISTER(utils.ClearBit(uint8(*reg), 6))
}

func (reg *PPU_STATUS_REGISTER) SetSpriteOverflow() {
	*reg = PPU_STATUS_REGISTER(utils.SetBit(uint8(*reg), 5))
}
func (reg *PPU_STATUS_REGISTER) ClearSpriteOverflow() {
	*reg = PPU_STATUS_REGISTER(utils.ClearBit(uint8(*reg), 5))
}
func (reg *PPU_STATUS_REGISTER) InVBlank() bool {
	return utils.HasBit(uint8(*reg), 7)
}

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
	if !utils.HasBit(bitsBefore, 7) && utils.HasBit(bitsAfter, 7) && ppu.Status.InVBlank() {
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

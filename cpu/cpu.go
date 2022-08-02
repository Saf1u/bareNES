package cpu

import (
	"emulator/ppu"
	"emulator/rom"
	"os"
)

const STACK uint8 = 0xfd
const (
	ACCUMULATOR                 = "Accumulator"
	RELATIVE                    = "relative"
	ABSOLUTE_INDIRECT           = "absindirect"
	IMPLIED                     = "implied"
	IMMEDIATE                   = "imm"
	ZERO_PAGE_X                 = "zpx"
	ABSOLUTE                    = "abs"
	ZERO_PAGE_Y                 = "zpy"
	ZERO_PAGE                   = "zpg"
	ABSOLUTE_X                  = "absx"
	ABSOLUTE_Y                  = "absy"
	INDIRECT_X                  = "indx"
	INDIRECT_Y                  = "indy"
	CARRY_FLAG                  = 0
	ZERO_FLAG                   = 1
	INTERRUPT_FLAG              = 2
	DECIMAL_FLAG                = 3
	BREAK_FLAG                  = 4
	OVERFLOW_FLAG               = 6
	NEGATIVE_FLAG               = 7
	INDIRECT                    = "ind"
	CPU_RAM_START        uint16 = 0x0000
	CPU_RAM_END          uint16 = 0x1fff
	PPU_DATA_REGISTER    uint16 = 0x2007
	PPU_CONTROL_REGISTER uint16 = 0x2000
	PPU_ADDRESS_REGISTER uint16 = 0x2006
	PPU_MASK_REGISTER    uint16 = 0x2001
	PPU_STATUS_REGISTER  uint16 = 0x2002
	PPU_OAM_ADDRESS      uint16 = 0x2003
	PPU_OAM_DATA         uint16 = 0x2004
	PPU_SCROLL_REGISTER  uint16 = 0x2005
	PPU_OAM_DMA          uint16 = 0x4014
	PPU_REGISTERS        uint16 = 0x2008
	PPU_REGISTERS_END    uint16 = 0x3fff
	PROG_ROM_START       uint16 = 0x8000
	PROG_ROM_END         uint16 = 0xffff
	CTRL_ONE                    = 0x4016
	CTRL_TWO                    = 0x4017
)

const STACK_PAGE uint16 = 0x0100

//Cpu composes of a 6502 register set and addressable memory
type Cpu struct {
	xRegister      uint8
	aRegister      uint8
	yRegister      uint8
	stackPtr       uint8
	pc             uint16
	statusRegister uint8
	CpuBus         Bus
	jumped         bool
}
type JoyPad struct {
	buttons  uint8
	strobe   bool
	bitIndex uint8
}

type Bus struct {
	cpuRam   [2048]uint8
	rom      *rom.Rom
	Ppu      *ppu.Ppu
	cpuTicks int
	Event    chan int
	done     chan int
	Pad      JoyPad
}

func (b *Bus) tick(amount int) {
	b.cpuTicks += amount

	res := b.Ppu.Tick(3 * amount)

	if res {

		b.Ppu.ShowTiles()

		b.Event <- 0
		<-b.done

	}
}

func (b *Bus) WriteSingleByte(addr uint16, data uint8) {
	switch {
	case addr >= CPU_RAM_START && addr <= CPU_RAM_END:
		addr = mirror(addr)
		b.cpuRam[addr] = data
	case addr == PPU_DATA_REGISTER:
		b.Ppu.WriteData(data)
	case addr == PPU_ADDRESS_REGISTER:
		b.Ppu.AddrRegister.Update(data)
	case addr == PPU_MASK_REGISTER:
		b.Ppu.Mask.Update(data)
	case addr == PPU_OAM_ADDRESS:
		b.Ppu.OamAddr.WriteAddressOam(data)
	case addr == PPU_SCROLL_REGISTER:
		b.Ppu.Scroll.Update(data)
	case addr == PPU_OAM_DATA:
		b.Ppu.WriteDataOam(data)
	case addr == PPU_OAM_DMA:
		start := uint16(data) << 8
		end := uint16(data)<<8 | 0x00FF
		oamData := []uint8{}
		for i := start; i <= end; i++ {
			oamData = append(oamData, b.ReadSingleByte(i))
		}
		b.Ppu.WriteOamDMA(oamData)
	case addr == PPU_CONTROL_REGISTER:
		b.Ppu.WriteToCtrl(data)
	case addr == CTRL_ONE:
		b.Pad.write(data)
	case addr == CTRL_TWO:
		b.Pad.write(data)
	case addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_END:
		addr = addr & 0b0010000000000111
		b.WriteSingleByte(addr, data)

	}

}
func (pad *JoyPad) write(data uint8) {
	//tf?
	if data&1 == 1 {

		pad.strobe = true
		pad.bitIndex = 0
		return
	}

	pad.strobe = false

}

func (pad *JoyPad) read() uint8 {

	if pad.bitIndex > 7 {
		return 1
	}
	current := (pad.buttons >> (pad.bitIndex)) & uint8(0b1)
	if !pad.strobe && pad.bitIndex <= 7 {
		pad.bitIndex++
	}
	return current

}
func (pad *JoyPad) Set(val bool, pos int) {
	if val {
		pad.buttons = setBit(pad.buttons, pos)
		return
	}
	pad.buttons = clearBit(pad.buttons, pos)
}

func (b *Bus) WriteDoubleByte(addr uint16, data uint16) {

	low := uint8(data & 0x00FF)
	hi := uint8((data) >> 8)
	b.WriteSingleByte(addr, low)
	b.WriteSingleByte(addr+1, hi)
}

func (b *Bus) ReadSingleByte(addr uint16) uint8 {
	var data uint8 = 0
	switch {
	case addr >= CPU_RAM_START && addr <= CPU_RAM_END:
		addr = mirror(addr)
		data = b.cpuRam[addr]
	case addr >= PROG_ROM_START && addr <= PROG_ROM_END:
		data = b.rom.ReadRom(addr)
	case addr == PPU_DATA_REGISTER:
		data = b.Ppu.ReadData()
	case addr == PPU_STATUS_REGISTER:
		data = b.Ppu.ReadStatus()
	case addr == PPU_OAM_DATA:
		data = b.Ppu.ReadDataOamRegister()
	case addr == CTRL_ONE:
		data = b.Pad.read()
	case addr == CTRL_TWO:
		data = b.Pad.read()
	case addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_END:
		addr = addr & 0b0010000000000111
		b.ReadSingleByte(addr)

	}
	return data

}
func (b *Bus) ReadDoubleByte(addr uint16) uint16 {

	var low uint16 = uint16(b.ReadSingleByte(addr))
	var hi uint16 = uint16(b.ReadSingleByte(addr + 1))
	res := (hi << 8) | low
	return res

}

func (c *Cpu) Acc() uint8 {
	return c.aRegister
}
func (c *Cpu) Stat() uint8 {
	return c.statusRegister
}
func (c *Cpu) addrMode(mode string) uint16 {
	var dataLocation uint16
	switch {
	case mode == RELATIVE:
		loc := c.pc + 1
		toJump := int8(c.CpuBus.ReadSingleByte(loc))
		dataLocation = uint16(toJump)
	case mode == IMMEDIATE:
		dataLocation = c.pc + 1
	case mode == ABSOLUTE:
		dataLocation = c.CpuBus.ReadDoubleByte(c.pc + 1)
	case mode == ZERO_PAGE_X:
		data := c.CpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.xRegister
		dataLocation = uint16(n)
	case mode == ZERO_PAGE_Y:
		data := c.CpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.yRegister
		dataLocation = uint16(n)
	case mode == ABSOLUTE_X:
		data := c.CpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
	case mode == ABSOLUTE_Y:
		data := c.CpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
	case mode == INDIRECT_X:
		base := uint16(c.CpuBus.ReadSingleByte(c.pc+1) + (c.xRegister))
		low := uint16(c.CpuBus.ReadSingleByte(uint16(base)))
		temp := uint8(base + 1)
		hi := uint16(c.CpuBus.ReadSingleByte(uint16(temp)))
		dataLocation = (hi << 8) | low
	case mode == INDIRECT_Y:
		pos := uint16(c.CpuBus.ReadSingleByte(c.pc + 1))
		low := c.CpuBus.ReadSingleByte(pos)
		temp := uint8(pos + 1)
		hi := c.CpuBus.ReadSingleByte(uint16(temp))
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)
	case mode == ZERO_PAGE:
		data := c.CpuBus.ReadSingleByte(c.pc + 1)
		dataLocation = uint16(data)
	case mode == INDIRECT:
		dataLocation = c.CpuBus.ReadDoubleByte(c.pc + 1)

	}

	return dataLocation
}

func (c *Cpu) Init(rom *rom.Rom, event chan int, done chan int) {

	c.CpuBus = Bus{}
	c.CpuBus.cpuTicks = 0
	c.CpuBus.rom = rom
	c.CpuBus.Ppu = ppu.NewPPU(rom.CharRom, rom.MirorType)
	c.CpuBus.Ppu.PpuTicks = c.CpuBus.cpuTicks * 3
	c.CpuBus.Event = event
	c.CpuBus.done = done
	c.xRegister = 0
	c.aRegister = 0
	c.yRegister = 0
	c.statusRegister = 0x24
	c.stackPtr = STACK
	c.pc = c.CpuBus.ReadDoubleByte(0xfffc)
}
func (c *Cpu) SEC() {
	c.statusRegister = (setBit(c.statusRegister, CARRY_FLAG))

}
func (c *Cpu) CLC() {
	c.statusRegister = (clearBit(c.statusRegister, CARRY_FLAG))
}
func (c *Cpu) GetBit(pos int) uint8 {
	return getBit(c.statusRegister, pos)
}
func (c *Cpu) SetZero() {

	c.statusRegister = (setBit(c.statusRegister, ZERO_FLAG))
}
func (c *Cpu) ClearZero() {
	c.statusRegister = (clearBit(c.statusRegister, ZERO_FLAG))
}
func (c *Cpu) SEI() {
	c.statusRegister = (setBit(c.statusRegister, INTERRUPT_FLAG))
}
func (c *Cpu) CLI() {
	c.statusRegister = (clearBit(c.statusRegister, INTERRUPT_FLAG))
}
func (c *Cpu) SED() {
	c.statusRegister = (setBit(c.statusRegister, DECIMAL_FLAG))
}
func (c *Cpu) CLD() {
	c.statusRegister = (clearBit(c.statusRegister, DECIMAL_FLAG))
}
func (c *Cpu) SetBreak() {
	c.statusRegister = (setBit(c.statusRegister, BREAK_FLAG))
}
func (c *Cpu) ClearBreak() {
	c.statusRegister = (clearBit(c.statusRegister, BREAK_FLAG))
}
func (c *Cpu) SetOverflow() {
	c.statusRegister = (setBit(c.statusRegister, OVERFLOW_FLAG))
}
func (c *Cpu) CLV() {
	c.statusRegister = (clearBit(c.statusRegister, OVERFLOW_FLAG))
}
func (c *Cpu) SetNegative() {
	c.statusRegister = (setBit(c.statusRegister, NEGATIVE_FLAG))
}
func (c *Cpu) ClearNegative() {
	c.statusRegister = (clearBit(c.statusRegister, NEGATIVE_FLAG))
}

func (c *Cpu) alterZeroAndNeg(data uint8) {
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(data, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}
func (c *Cpu) ISC(mode string) {
	c.INC(mode)
	c.SBC(mode)
}
func (c *Cpu) LAX(mode string) {
	c.LDA(mode)
	c.xRegister = c.aRegister
}

func (c *Cpu) LDX(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) LDY(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.yRegister = data
}
func (c *Cpu) SBC(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)

	data = ^data

	if hasBit(c.statusRegister, CARRY_FLAG) {
		data++
	}
	t := (c.aRegister) + (data)
	temp := uint8(t)

	if (!hasBit(c.aRegister, 7) && !hasBit(data, 7)) && (hasBit(temp, 7)) {
		c.SetOverflow()
		if (uint16(c.aRegister) + uint16(data)) != uint16(temp) {
			c.SEC()
		} else {
			c.CLC()
		}
	} else {
		if (hasBit(c.aRegister, 7) && hasBit(data, 7)) && (!hasBit(temp, 7)) {

			c.SetOverflow()
			if (uint16(c.aRegister) + uint16(data)) != uint16(temp) {
				c.SEC()
			} else {
				c.CLC()
			}

		} else {

			c.CLV()

			if (!hasBit(c.aRegister, 7) && hasBit(data, 7)) || (hasBit(c.aRegister, 7) && !hasBit(data, 7)) {

				if (uint16(c.aRegister) + uint16(data)) != uint16(temp) {
					c.SEC()
				} else {
					if data != 0 {
						c.CLC()
					}
				}

			}
		}

	}
	c.aRegister = temp
	if c.aRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.aRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) ADC(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)

	t := (c.aRegister) + (data)
	if hasBit(c.statusRegister, CARRY_FLAG) {
		t++
	}
	c.CLC()
	temp := uint8(t)

	if (!hasBit(c.aRegister, 7) && !hasBit(data, 7)) && (hasBit(temp, 7)) {
		c.SetOverflow()
	} else {
		if (hasBit(c.aRegister, 7) && hasBit(data, 7)) && (!hasBit(temp, 7)) {
			c.SetOverflow()
		} else {
			c.CLV()
			if (!hasBit(c.aRegister, 7) && hasBit(data, 7)) || (hasBit(c.aRegister, 7) && !hasBit(data, 7)) {
				if (uint16(c.aRegister) + uint16(data)) != uint16(temp) {
					c.SEC()
				} else {
					c.CLC()
				}
			}
		}

	}
	c.aRegister = temp
	if c.aRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.aRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}
}

func (c *Cpu) STA(mode string) {
	loc := c.addrMode(mode)
	c.CpuBus.WriteSingleByte(loc, c.aRegister)
}
func (c *Cpu) STX(mode string) {
	loc := c.addrMode(mode)

	c.CpuBus.WriteSingleByte(loc, c.xRegister)
}
func (c *Cpu) STY(mode string) {
	loc := c.addrMode(mode)

	c.CpuBus.WriteSingleByte(loc, c.yRegister)
}
func (c *Cpu) SAX(mode string) {
	loc := c.addrMode(mode)
	data := c.aRegister & c.xRegister
	c.CpuBus.WriteSingleByte(loc, data)
}

func (c *Cpu) TAX() {
	data := c.aRegister
	c.alterZeroAndNeg(data)
	c.xRegister = data
}
func (c *Cpu) TAY() {
	data := c.aRegister
	c.alterZeroAndNeg(data)
	c.yRegister = data
}

func (c *Cpu) TXA() {
	data := c.xRegister
	c.alterZeroAndNeg(data)
	c.aRegister = data
}
func (c *Cpu) TYA() {
	data := c.yRegister
	c.alterZeroAndNeg(data)

	c.aRegister = data
}
func (c *Cpu) TXS() {
	data := c.xRegister
	c.stackPtr = data
}

func (c *Cpu) TSX() {
	data := c.stackPtr
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) AND(mode string, hidden ...*uint8) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	if len(hidden) != 0 {
		data = *hidden[0]
	}
	c.aRegister = data & c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) ORA(mode string, hidden ...*uint8) {
	loc := c.addrMode(mode)

	data := c.CpuBus.ReadSingleByte(loc)
	if len(hidden) != 0 {
		data = *(hidden[0])
	}
	c.aRegister = data | c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) EOR(mode string, hidden ...*uint8) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	if len(hidden) != 0 {
		data = *hidden[0]
	}
	c.aRegister = data ^ c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) INX() {
	c.xRegister++
	c.alterZeroAndNeg(c.xRegister)

}
func (c *Cpu) INY() {
	c.yRegister++
	c.alterZeroAndNeg(c.yRegister)

}

func (c *Cpu) DEX() {
	c.xRegister--
	c.alterZeroAndNeg(c.xRegister)

}

func (c *Cpu) DEY() {
	c.yRegister--
	c.alterZeroAndNeg(c.yRegister)

}

func (c *Cpu) INC(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	data++
	c.CpuBus.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)

}

func (c *Cpu) DEC(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	data--
	c.CpuBus.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)

}

func (c *Cpu) CMP(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	temp := c.aRegister - data
	if c.aRegister >= data {
		c.SEC()
	} else {
		c.CLC()
	}
	if c.aRegister == data {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(temp, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) CPX(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	temp := c.xRegister - data
	if c.xRegister >= data {
		c.SEC()
	} else {
		c.CLC()
	}
	if c.xRegister == data {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(temp, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}
func (c *Cpu) CPY(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	temp := c.yRegister - data
	if c.yRegister >= data {
		c.SEC()
	} else {
		c.CLC()
	}
	if c.yRegister == data {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(temp, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) BIT(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	temp := data & c.aRegister
	if temp == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(data, 6) {
		c.SetOverflow()
	} else {
		c.CLV()
	}
	if hasBit(data, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) LSR(mode string) {
	var data uint8
	var loc uint16
	if mode == ACCUMULATOR {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.CpuBus.ReadSingleByte(loc)
	}
	if hasBit(data, 0) {
		c.SEC()
	} else {
		c.CLC()
	}
	data = data >> 1
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.CpuBus.WriteSingleByte(loc, data)
	}
	c.ClearNegative()
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}

}
func (c *Cpu) SLO(mode string) {
	c.ASL(mode)
	c.ORA(mode)
}

func (c *Cpu) ASL(mode string) {
	var data uint8
	var loc uint16
	if mode == ACCUMULATOR {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.CpuBus.ReadSingleByte(loc)
	}
	if hasBit(data, 7) {
		c.SEC()
	} else {
		c.CLC()
	}
	data = data << 1
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.CpuBus.WriteSingleByte(loc, data)
	}
	c.alterZeroAndNeg(data)

}

func (c *Cpu) ROL(mode string) {
	var data uint8
	var loc uint16
	if mode == ACCUMULATOR {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.CpuBus.ReadSingleByte(loc)
	}

	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 7)
	data = data << 1
	if temp > 0 {
		data = setBit(data, 0)
	} else {
		data = clearBit(data, 0)
	}
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.CpuBus.WriteSingleByte(loc, data)
	}

	if templast > 0 {
		c.SEC()
	} else {
		c.CLC()
	}
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	templast = getBit(data, 7)

	if templast > 0 {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}
func (c *Cpu) SRE(mode string) {

	c.LSR(mode)
	c.EOR(mode)

}
func (c *Cpu) AXS(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	c.xRegister = c.xRegister & c.aRegister
	if data <= c.xRegister {
		c.SEC()
	}
	c.xRegister = c.xRegister - data

	c.alterZeroAndNeg(c.xRegister)
}

func (c *Cpu) ARR(mode string) {
	loc := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(loc)
	c.AND(mode, &data)
	c.aRegister = data & c.aRegister
	c.ROR(ACCUMULATOR)
	temp := c.Acc()
	if hasBit(temp, 5) && hasBit(temp, 6) {
		c.SEC()
		c.CLV()
	}
	if !hasBit(temp, 5) && !hasBit(temp, 6) {
		c.CLC()
		c.CLV()
	}
	if hasBit(temp, 5) && !hasBit(temp, 6) {
		c.CLC()
		c.SetOverflow()
	}
	if !hasBit(temp, 5) && hasBit(temp, 6) {
		c.SetOverflow()
		c.SEC()
	}

	c.alterZeroAndNeg(temp)

}

func (c *Cpu) ROR(mode string) {
	var data uint8
	var loc uint16
	if mode == ACCUMULATOR {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.CpuBus.ReadSingleByte(loc)
	}
	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 0)
	data = data >> 1
	if temp > 0 {
		data = setBit(data, 7)
	} else {
		data = clearBit(data, 7)
	}
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.CpuBus.WriteSingleByte(loc, data)
	}

	if templast > 0 {
		c.SEC()
	} else {
		c.CLC()
	}
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}

	if temp > 0 {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) JMP(mode string) {

	if mode == ABSOLUTE {
		c.pc = c.addrMode(ABSOLUTE)
	} else {
		loc := c.CpuBus.ReadDoubleByte(c.pc + 1)
		//6502 HAS A WEIRD WRAPAROUND BUG THAT CAUSES AN ADDRESS TO BE READ BACKWARD IN AN INDIRECT JUMP WE NEED TO REMAIN TRUE TO THIS
		//
		if loc&0x00ff == 0x00ff {
			low := uint16(c.CpuBus.ReadSingleByte(loc))
			hi := uint16(c.CpuBus.ReadSingleByte(loc & 0xFF00))
			fin := hi<<8 | low
			c.pc = fin
		} else {
			fin := c.CpuBus.ReadDoubleByte(loc)
			c.pc = fin
		}

	}

}

func (c *Cpu) BMI(tick int) {
	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 7) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) BPL(tick int) {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 7) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}

func (c *Cpu) BVS(tick int) {

	toJump := c.addrMode(RELATIVE)

	if hasBit(c.statusRegister, 6) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}

func (c *Cpu) BVC(tick int) {

	//location of perand to jump too in mem not acc value itself is loc
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 6) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + (toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) BCC(tick int) {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 0) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {

		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) BEQ(tick int) {

	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 1) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) BCS(tick int) {
	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 0) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) BNE(tick int) {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 1) {
		c.jumped = true
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	} else {
		if tick == 4 {
			tick -= 2
		} else {
			tick--
		}
	}
	c.CpuBus.tick(tick)
}
func (c *Cpu) ANC(mode string) {

	c.AND(mode)
	if hasBit(c.aRegister, NEGATIVE_FLAG) {
		c.SEC()
	} else {
		c.CLC()
	}
}

func (c *Cpu) ALR(mode string) {

	c.AND(mode)
	c.LSR(ACCUMULATOR)
}

func (c *Cpu) RRA(mode string) {

	c.ROR(mode)
	c.ADC(mode)
}

func (c *Cpu) DCP(mode string) {
	location := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(location)
	data--
	c.CpuBus.WriteSingleByte(location, data)
	if data <= c.aRegister {
		c.SEC()
	}
	data = c.aRegister - data
	c.alterZeroAndNeg(data)
}

func (c *Cpu) XAA(mode string) {
	location := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(location)
	c.aRegister = c.xRegister
	c.alterZeroAndNeg(c.aRegister)
	c.AND(mode, &data)

}

func (c *Cpu) LAS(mode string) {
	location := c.addrMode(mode)
	data := c.CpuBus.ReadSingleByte(location)
	c.stackPtr = c.stackPtr & data
	c.aRegister = c.stackPtr & data
	c.xRegister = c.stackPtr & data
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) TAS(mode string) {
	data := c.aRegister & c.xRegister
	c.stackPtr = data
	loc := c.addrMode(mode)
	data = (uint8(loc>>8) + 1) & data
	c.CpuBus.WriteSingleByte(loc, data)

}
func (c *Cpu) AHX(mode string) {
	data := c.aRegister & c.xRegister
	loc := c.addrMode(mode)
	data = (uint8(loc>>8) + 1) & data
	c.CpuBus.WriteSingleByte(loc, data)
}
func (c *Cpu) SHX(mode string) {
	data := c.xRegister
	loc := c.addrMode(mode)
	data = (uint8(loc>>8) + 1) & data
	c.CpuBus.WriteSingleByte(loc, data)
}
func (c *Cpu) SHY(mode string) {
	data := c.yRegister
	loc := c.addrMode(mode)
	data = (uint8(loc>>8) + 1) & data
	c.CpuBus.WriteSingleByte(loc, data)
}

func (c *Cpu) RLA(mode string) {

	c.ROL(mode)
	c.AND(mode)

}

func (c *Cpu) JSR() {
	c.jumped = true
	//we need to make sure we increment within the same cycle
	addr := c.addrMode(ABSOLUTE)
	c.PushDouble(c.pc + 2)
	c.pc = addr
}
func (c *Cpu) RTS() {
	c.jumped = true
	val := c.PopDouble()
	c.pc = val + 1

}
func (c *Cpu) PHA() {
	acc := c.Acc()
	c.Push(acc)
}
func (c *Cpu) PHP() {
	reg := c.statusRegister
	reg = setBit(reg, BREAK_FLAG)
	reg = setBit(reg, 5)
	c.Push(reg)
}
func (c *Cpu) PLA() {
	acc := c.Pop()
	c.alterZeroAndNeg(acc)
	c.aRegister = acc
}
func (c *Cpu) ExecuteNMI() {
	temp := setBit(c.statusRegister, 5)
	temp = clearBit(temp, BREAK_FLAG)
	//hardware interrupt by ppu clear break??
	c.PushDouble(c.pc)
	c.Push(temp)
	c.SEI()
	data := c.CpuBus.ReadDoubleByte(0xfffa)
	c.pc = data
	c.CpuBus.tick(2)
}
func (c *Cpu) RTI() {
	c.jumped = true
	c.statusRegister = c.Pop()
	c.pc = c.PopDouble()
	c.ClearBreak()
	c.statusRegister = setBit(c.statusRegister, 5)
}
func (c *Cpu) BRK() {
	c.Push(c.statusRegister)
	c.PushDouble(c.pc)
	c.SEI()
	c.pc = c.CpuBus.ReadDoubleByte(0xfffe)
}

func (c *Cpu) PLP() {
	reg := c.Pop()
	c.statusRegister = reg
	c.ClearBreak()
	c.statusRegister = setBit(c.statusRegister, 5)
}

func (c *Cpu) LDA(mode string) {
	loc := c.addrMode(mode)

	data := c.CpuBus.ReadSingleByte(loc)
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(data, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

	c.aRegister = data

}

func (c *Cpu) NOP(mode string) {
	if mode == IMPLIED {
		return
	}

}

func (c *Cpu) Push(val uint8) {
	loc := 0x0100 + uint16(c.stackPtr)
	c.CpuBus.WriteSingleByte(loc, val)
	c.stackPtr--
}
func (c *Cpu) PushDouble(val uint16) {
	hi := uint8(val >> 8)
	lo := uint8(val & 0x00FF)
	c.Push(hi)
	c.Push(lo)
}
func (c *Cpu) PopDouble() uint16 {
	lo := uint16(c.Pop())
	hi := uint16(c.Pop())
	val := (hi<<8 | lo)
	return val
}
func (c *Cpu) Pop() uint8 {
	c.stackPtr++
	//stack grows down
	loc := 0x0100 + uint16(c.stackPtr)
	temp := c.CpuBus.ReadSingleByte(loc)
	return temp
}

func (c *Cpu) Run() {
	fl, err := os.OpenFile("tracer.log", os.O_RDWR, 0755)
	defer fl.Close()
	if err != nil {
		return
	}
	for {
		c.jumped = false

		if c.CpuBus.Ppu.PollNmi() {
			c.CpuBus.Ppu.NmiOcurred = false
			c.ExecuteNMI()
		}

		//temp := c.pc

		location := c.CpuBus.ReadSingleByte(c.pc)
		//fmt.Printf("%x %x %x %x\n", location, c.CpuBus.ReadDoubleByte(c.pc+1), c.pc, c.statusRegister)
		mode := getAddrMode(location)
		tick := c.TraceExecution(mode, fl)

		switch location {
		case 0x00:
			//c.BRK()
			return
		case 0x10:
			c.BPL(tick)
		case 0x20:
			c.JSR()
		case 0x30:
			c.BMI(tick)
		case 0x40:
			c.RTI()
		case 0x50:
			c.BVC(tick)
		case 0x60:
			c.RTS()
		case 0x70:
			c.BVS(tick)
		case 0x90:
			c.BCC(tick)
		case 0xA0, 0xA4, 0xB4, 0xAC, 0xBC:
			c.LDY(mode)
		case 0xB0:
			c.BCS(tick)

		case 0xD0:
			c.BNE(tick)

		case 0xF0:
			c.BEQ(tick)
		case 0x01, 0x09, 0x05, 0x15, 0x0D, 0x1D, 0x19, 0x11:
			c.ORA(mode)
		case 0x21, 0x29, 0x25, 0x35, 0x2D, 0x3D, 0x39, 0x31:
			c.AND(mode)

		case 0xd1, 0xC9, 0xC5, 0xD5, 0xCD, 0xDD, 0xD9, 0xC1:
			c.CMP(mode)

		case 0x24:
			c.BIT(mode)
		case 0x84, 0x94, 0x8C:
			c.STY(mode)

		case 0xC4, 0xC0, 0xCC:
			c.CPY(mode)
		case 0xe4, 0xE0, 0xEC:
			c.CPX(mode)

		case 0x65, 0x69, 0x75, 0x6D, 0x7D, 0x79, 0x61, 0x71:
			c.ADC(mode)

		case 0xA5, 0xA9, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1:
			c.LDA(mode)
		case 0x4b:
			c.ALR(mode)
		case 0xf5, 0xE9, 0xE5, 0xED, 0xFD, 0xF9, 0xE1, 0xF1, 0xEB:
			c.SBC(mode)
		case 0x06, 0x0A, 0x16, 0x0E, 0x1E:
			c.ASL(mode)

		case 0x26:
			c.ROL(mode)
		case 0x36:
			c.ROL(mode)
		case 0x46, 0x4A, 0x56, 0x4E, 0x5E:
			c.LSR(mode)

		case 0x76, 0x6A, 0x66, 0x6E, 0x7E:
			c.ROR(mode)
		case 0x86, 0x96, 0x8E:
			c.STX(mode)
		case 0xA6, 0xA2, 0xB6, 0xAE, 0xBE:
			c.LDX(mode)

		case 0xc6, 0xD6, 0xCE, 0xDE:
			c.DEC(mode)
		case 0xA7, 0xB7, 0xAF, 0xBf, 0xA3, 0xB3, 0xab:
			c.LAX(mode)
		case 0x67, 0x77, 0x6F, 0x7f, 0x7B, 0x63, 0x73:
			c.RRA(mode)
		case 0xE7, 0xF7, 0xEF, 0xFf, 0xFB, 0xE3, 0xF3:
			c.ISC(mode)
		case 0x0b, 0x2b:
			c.ANC(mode)

		case 0x08:
			c.PHP()
		case 0x18:
			c.CLC()
		case 0x28:
			c.PLP()
		case 0x38:
			c.SEC()
		case 0x48:
			c.PHA()
		case 0x58:
			c.CLI()
		case 0x68:
			c.PLA()
		case 0x78:
			c.SEI()
		case 0x88:
			c.DEY()

		case 0x87, 0x97, 0x8f, 0x83:
			c.SAX(mode)
		case 0x8B:
			c.XAA(mode)
		case 0x98:
			c.TYA()
		case 0xA8:
			c.TAY()
		case 0xB8:
			c.CLV()
		case 0xc8:
			c.INY()
		case 0xD8:
			c.CLD()
			//lmao no decimal mode
		case 0xE8:
			c.INX()
		case 0xc7, 0xd7, 0xcf, 0xdf, 0xdb, 0xc3, 0xd3:
			c.DCP(mode)
		case 0xF8:
			c.SED()
		case 0x27, 0x37, 0x3f, 0x2f, 0x3b, 0x33, 0x23:
			c.RLA(mode)
			//lmao no decimal mode
		case 0x49, 0x45, 0x55, 0x4D, 0x5D, 0x59, 0x41, 0x51:
			c.EOR(mode)
		case 0x07, 0x17, 0x0f, 0x1f, 0x1b, 0x03, 0x13:
			c.SLO(mode)
		case 0x47, 0x57, 0x4f, 0x5f, 0x5b, 0x43, 0x53:
			c.SRE(mode)
		case 0xCB:
			c.AXS(mode)
		case 0x6b:
			c.ARR(mode)
		case 0x99, 0x85, 0x95, 0x8D, 0x9D, 0x81, 0x91:
			c.STA(mode)

		case 0x2a:
			c.ROL(mode)

		case 0x8a:
			c.TXA()
		case 0x9a:
			c.TXS()
		case 0xAa:
			c.TAX()
		case 0xba:
			c.TSX()
		case 0xca:
			c.DEX()
		case 0xea, 0x1A, 0x3A, 0x5A, 0x7A, 0xDA, 0xFA, 0x80, 0x82, 0x89, 0xC2, 0xE2, 0x04, 0x44, 0x64, 0x14, 0x34, 0x54, 0x74, 0xD4, 0xF4, 0x0C, 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
			c.NOP(mode)
		case 0x2c:
			c.BIT(mode)
		case 0x4c, 0x6C:
			c.jumped = true
			c.JMP(mode)

		case 0x2e:
			c.ROL(mode)
		case 0x3e:
			c.ROL(mode)

		case 0xee, 0xE6, 0xF6, 0xFE:
			c.INC(mode)
		case 0xbb:
			c.LAS(mode)
		case 0x9b:
			c.TAS(mode)
		case 0x9f, 0x93:
			c.AHX(mode)
		case 0x9E:
			c.SHX(mode)
		case 0x9c:
			c.SHY(mode)

		}
		switch location {
		case 0x10, 0x30, 0x50, 0x70, 0x90, 0xb0, 0xd0, 0xf0:
		default:
			c.CpuBus.tick(tick)
		}

		if !c.jumped {
			val := getNumber(location)
			if val == -1 {
				break
			}
			if c.pc+uint16(val) < c.pc {
				break
			}
			c.pc = c.pc + uint16(val)

		}
	}

}

package cpu

const programLocation = 0x8000

//Cpu composes of a 6502 register set and addressable memory
type Cpu struct {
	xRegister      uint8
	aRegister      uint8
	yRegister      uint8
	stackPtr       uint8
	pc             uint16
	statusRegister uint8
	mem            [0xFFFF]uint8
}

func (c *Cpu) Acc() uint8 {
	return c.aRegister
}
func (c *Cpu) Stat() uint8 {
	return c.statusRegister
}

var programLength int

//num of bytes to move pc depending on instruction
var pcIncrement = map[uint8]int{
	0x69: 2, 0x65: 2, 0x75: 2, 0x6D: 3, 0x7D: 3, 0x79: 3, 0x61: 2, 0x71: 2,
	0x29: 2, 0x25: 2, 0x35: 2, 0x2D: 3, 0x3D: 3, 0x39: 3, 0x21: 2, 0x31: 2,
	0x0A: 1, 0x06: 2, 0x16: 2, 0x0E: 3, 0x1E: 3,
	0x90: 2,
	0xb0: 2,
	0xF0: 2,
	0x24: 2, 0x2C: 3,
	0x30: 2,
	0xD0: 2,
	0x10: 2,
	0x00: 1,
	0x50: 2,
	0x70: 2,
	0x18: 1,
	0xD8: 1,
	0x58: 1,
	0xB8: 1,
	0xC9: 2, 0xC5: 2, 0xD5: 2, 0xCD: 3, 0xDD: 3, 0xD9: 3, 0xC1: 2, 0xD1: 2,
	0xE0: 2, 0xE4: 2, 0xEC: 3,
	0xC0: 2, 0xC4: 2, 0xCC: 3,
	0xC6: 2, 0xD6: 2, 0xCE: 3, 0xDE: 3,
	0xCA: 1,
	0x88: 1,
	0x49: 2, 0x45: 2, 0x55: 2, 0x4D: 3, 0x5D: 3, 0x59: 3, 0x41: 2, 0x51: 2,
	0xE6: 2, 0xF6: 2, 0xEE: 3, 0xFE: 3,
	0xE8: 1,
	0xC8: 1,
	0x4C: 3, 0x6C: 3,
	0x20: 3,
	0xA9: 2, 0xA5: 2, 0xB5: 2, 0xAD: 3, 0xBD: 3, 0xB9: 3, 0xA1: 2, 0xB1: 2,
	0xA2: 2, 0xA6: 2, 0xB6: 2, 0xAE: 3, 0xBE: 3,
	0xA0: 2, 0xA4: 2, 0xB4: 2, 0xAC: 3, 0xBC: 3,
	0x4A: 1, 0x46: 2, 0x56: 2, 0x4E: 3, 0x5E: 3,
	0xEA: 1,
	0x09: 2, 0x05: 2, 0x15: 2, 0x0D: 3, 0x1D: 3, 0x19: 3, 0x01: 2, 0x11: 2,
	0x48: 1,
	0x08: 1,
	0x68: 1,
	0x28: 1,
	0x2A: 1, 0x26: 2, 0x36: 2, 0x2E: 3, 0x3E: 3,
	0x6A: 1, 0x66: 2, 0x76: 2, 0x6E: 3, 0x7E: 3,
	0x40: 1,
	0x60: 1,
	0xE9: 2, 0xE5: 2, 0xF5: 2, 0xED: 3, 0xFD: 3, 0xF9: 3, 0xE1: 2, 0xF1: 2,
	0x38: 1,
	0xF8: 1,
	0x78: 1,
	0x85: 2, 0x95: 2, 0x8D: 3, 0x9D: 3, 0x99: 3, 0x81: 2, 0x91: 2,
	0x86: 2, 0x96: 2, 0x8E: 3,
	0x84: 2, 0x94: 2, 0x8C: 3,
	0xAA: 1,
	0xA8: 1,
	0xBA: 1,
	0x8A: 1,
	0x9A: 1,
	0x98: 1,
}

func (c *Cpu) ReadDoubleByte(addr uint16) uint16 {
	var low uint16 = uint16(c.mem[addr])
	var hi uint16 = uint16(c.mem[addr+1])
	res := (hi << 8) | low
	return res
}
func (c *Cpu) addrMode(mode string) uint16 {
	var dataLocation uint16
	switch {
	case mode == "imm":
		dataLocation = c.pc
		break
	case mode == "zp":
		dataLocation = uint16(c.ReadSingleByte(c.pc))
	case mode == "abs":
		dataLocation = c.ReadDoubleByte(c.pc)
	case mode == "zpx":
		data := c.ReadSingleByte(c.pc)
		var c uint8 = data + c.xRegister
		dataLocation = uint16(c)
	case mode == "zpy":
		data := c.ReadSingleByte(c.pc)
		var c uint8 = data + c.yRegister
		dataLocation = uint16(c)
	case mode == "absx":
		data := c.ReadDoubleByte(c.pc)
		dataLocation = data + uint16(c.xRegister)
	case mode == "absy":
		data := c.ReadDoubleByte(c.pc)
		dataLocation = data + uint16(c.yRegister)
	case mode == "indx":
		base := c.ReadSingleByte(c.pc) + c.xRegister
		low := uint16(c.ReadSingleByte(uint16(base)))
		hi := uint16(c.ReadSingleByte(uint16(base) + 1))
		dataLocation = (hi << 8) | low
	case mode == "indy":
		pos := uint16(c.ReadSingleByte(c.pc))
		low := c.ReadSingleByte(pos)
		hi := c.ReadSingleByte(pos + 1)
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)

	}
	return dataLocation
}

func (c *Cpu) set() {
	c.xRegister = 0
	c.aRegister = 0
	c.yRegister = 0
	c.pc = c.ReadDoubleByte(0xfffc)
	c.statusRegister = 0
}

func (c *Cpu) LDA(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.aRegister = data
}

func (c *Cpu) LDX(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.xRegister = data
}

func (c *Cpu) LDY(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.xRegister = data
}

func (c *Cpu) STA(mode string) {
	loc := c.addrMode(mode)
	c.WriteSingleByte(loc, c.aRegister)
}
func (c *Cpu) STX(mode string) {
	loc := c.addrMode(mode)
	c.WriteSingleByte(loc, c.xRegister)
}
func (c *Cpu) STY(mode string) {
	loc := c.addrMode(mode)
	c.WriteSingleByte(loc, c.yRegister)
}
func (c *Cpu) TAX() {
	data := c.aRegister
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.xRegister = data
}
func (c *Cpu) TAY() {
	data := c.aRegister
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.yRegister = data
}

func (c *Cpu) TXA() {
	data := c.xRegister
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.aRegister = data
}
func (c *Cpu) TYA() {
	data := c.yRegister
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.aRegister = data
}
func (c *Cpu) TXS() {
	data := c.xRegister
	c.stackPtr = data
}

func (c *Cpu) TSX() {
	data := c.stackPtr
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

	c.xRegister = data
}

func (c *Cpu) AND(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data & c.aRegister
	if c.aRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.aRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) ORA(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data | c.aRegister
	if c.aRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.aRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) EOR(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data ^ c.aRegister
	if c.aRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.aRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) INX() {
	c.xRegister++
	if c.xRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.xRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}
func (c *Cpu) INY() {
	c.yRegister++
	if c.yRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.yRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) DEX() {
	c.xRegister--
	if c.xRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.xRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) DEY() {
	c.yRegister--
	if c.yRegister == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(c.yRegister, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) INC(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	data++
	c.WriteSingleByte(loc, data)
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) DEC(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	data--
	c.WriteSingleByte(loc, data)
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) CMP(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	temp := c.aRegister - data
	if c.aRegister >= data {
		c.statusRegister = (setBit(c.statusRegister, 0))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 0))
	}
	if c.aRegister == data {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(temp, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) CPX(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	temp := c.xRegister - data
	if c.xRegister >= data {
		c.statusRegister = (setBit(c.statusRegister, 0))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 0))
	}
	if c.xRegister == data {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(temp, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}
func (c *Cpu) CPY(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	temp := c.yRegister - data
	if c.yRegister >= data {
		c.statusRegister = (setBit(c.statusRegister, 0))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 0))
	}
	if c.yRegister == data {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(temp, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) BIT(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	temp := data & c.aRegister
	if temp == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}
	if hasBit(data, 6) {
		c.statusRegister = (setBit(c.statusRegister, 6))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 6))
	}
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}

}

func (c *Cpu) LSR(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	if hasBit(data, 0) {
		c.statusRegister = (setBit(c.statusRegister, 0))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 0))
	}
	data = data >> 1
	c.WriteSingleByte(loc, data)
	c.statusRegister = (clearBit(c.statusRegister, 7))
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}

}

func (c *Cpu) ASL(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 0))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 0))
	}
	data = data << 1
	c.WriteSingleByte(loc, data)
	if hasBit(data, 7) {
		c.statusRegister = (setBit(c.statusRegister, 7))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 7))
	}
	if data == 0 {
		c.statusRegister = (setBit(c.statusRegister, 1))
	} else {
		c.statusRegister = (clearBit(c.statusRegister, 1))
	}

}

// func (c *Cpu) ROL(mode string) {
// 	loc := c.addrMode(mode)
// 	data := c.ReadSingleByte(loc)
// 	oldcarry
// 	if hasBit(c.statusRegister, 0) {

// 	}

// }
// func (c *Cpu) ROR(mode string) {
// 	loc := c.addrMode(mode)
// 	data := c.ReadSingleByte(loc)
// 	oldcarry
// 	if hasBit(c.statusRegister, 0) {

// 	}

// }
func (c *Cpu) JMP(mode string) {
	loc := c.addrMode("abs")
	if mode == "absind" {
		c.pc = loc
	} else {
		//6502 HAS A WEIRD WRAPAROUND BUG THAT CAUSES AN ADDRESS TO BE READ BACKWARD IN AN INDIRECT JUMP WE NEED TO REMAIN TRUE TO THIS
		//
		if loc&0x00ff == 0x00ff {
			low := uint16(c.ReadSingleByte(loc))
			hi := uint16(c.ReadSingleByte(loc & 0xFF00))
			fin := hi<<8 | low
			c.pc = fin
		} else {
			c.pc = c.ReadDoubleByte(loc)
		}
	}

}

// func (c *Cpu) BMI() {
// 	loc := c.addrMode("imm")
// 	if hasBit(c.statusRegister, 7) {
// 		c.pc = c.pc + loc

// 	}
// }

func (c *Cpu) run() {

	for {
		inst := c.mem[c.pc]
		c.pc++
		switch {
		case inst == 0x00:
			return
		case inst == 0xa9:
			c.LDA("imm")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xA5:
			c.LDA("zp")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xb5:
			c.LDA("zpx")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xad:
			c.LDA("abs")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xbd:
			c.LDA("absx")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xb9:
			c.LDA("absy")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xa1:
			c.LDA("indx")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc
		case inst == 0xb1:
			c.LDA("indy")
			inc := uint16(pcIncrement[inst]) - 1
			c.pc = c.pc + inc

		}
		if c.pc == programLocation {
			break
		}
	}
}

//WriteSingleByte writes single byte to mem
func (c *Cpu) WriteSingleByte(addr uint16, data uint8) {

	c.mem[addr] = data

}

//ReadSingleByte writes single byte to mem
func (c *Cpu) ReadSingleByte(addr uint16) uint8 {
	data := c.mem[addr]
	return data
}

func (c *Cpu) WriteDoubleByte(addr uint16, data uint16) {
	low := uint8(data & 0x00FF)
	c.mem[addr] = low
	hi := uint8((data) >> 8)
	c.mem[addr+1] = hi

}

func (c *Cpu) LoadToRom(data []uint8) {
	copy(c.mem[programLocation:programLocation+len(data)], data)
	c.WriteDoubleByte(0xFFFC, programLocation)
}

func (c *Cpu) LoadToRomandStart(data []uint8) {
	programLength = programLocation + len(data)
	c.LoadToRom(data)
	c.set()
	c.run()
}

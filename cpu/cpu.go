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

const STACK uint16 = 0

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

const (
	ACCUMULATOR    = "Accumulator"
	IMMEDIATE      = "imm"
	ZERO_PAGE_X    = "zpx"
	ABSOLUTE       = "abs"
	ZERO_PAGE_Y    = "zpy"
	ABSOLUTE_X     = "absx"
	ABSOLUTE_Y     = "absy"
	INDIRECT_X     = "indx"
	INDRECT_Y      = "indy"
	CARRY_FLAG     = 0
	ZERO_FLAG      = 1
	INTERRUPT_FLAG = 2
	BREAK_FLAG     = 4
	OVERFLOW_FLAG  = 6
	NEGATIVE_FLAG  = 7
)
const STACK_PAGE uint16 = 0x0100

func (c *Cpu) Acc() uint8 {
	return c.aRegister
}
func (c *Cpu) Stat() uint8 {
	return c.statusRegister
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
	case mode == IMMEDIATE:
		dataLocation = c.pc
	case mode == ZERO_PAGE_X:
		dataLocation = uint16(c.ReadSingleByte(c.pc))
	case mode == ABSOLUTE_X:
		dataLocation = c.ReadDoubleByte(c.pc)
	case mode == ZERO_PAGE_X:
		data := c.ReadSingleByte(c.pc)
		var c uint8 = data + c.xRegister
		dataLocation = uint16(c)
	case mode == ZERO_PAGE_Y:
		data := c.ReadSingleByte(c.pc)
		var c uint8 = data + c.yRegister
		dataLocation = uint16(c)
	case mode == ABSOLUTE_X:
		data := c.ReadDoubleByte(c.pc)
		dataLocation = data + uint16(c.xRegister)
	case mode == ABSOLUTE_Y:
		data := c.ReadDoubleByte(c.pc)
		dataLocation = data + uint16(c.yRegister)
	case mode == INDIRECT_X:
		base := c.ReadSingleByte(c.pc) + c.xRegister
		low := uint16(c.ReadSingleByte(uint16(base)))
		hi := uint16(c.ReadSingleByte(uint16(base) + 1))
		dataLocation = (hi << 8) | low
	case mode == INDRECT_Y:
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

func (c *Cpu) LDX(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) LDY(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.xRegister = data
}
func (c *Cpu) ADC(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = c.aRegister + data
	if hasBit(c.statusRegister, CARRY_FLAG) {
		c.aRegister++
	}
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
func (c *Cpu) stackIncrement() uint16 {
	return uint16(c.stackPtr) + STACK_PAGE
}

func (c *Cpu) TSX() {
	data := c.stackPtr
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) AND(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data & c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) ORA(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data | c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) EOR(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = data ^ c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) INX() {
	c.xRegister++
	if c.xRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.xRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}
func (c *Cpu) INY() {
	c.yRegister++
	if c.yRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.yRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) DEX() {
	c.xRegister--
	if c.xRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.xRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) DEY() {
	c.yRegister--
	if c.yRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.yRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

}

func (c *Cpu) INC(mode string, hidden ...*uint8) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	data++
	c.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) DEC(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	data--
	c.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)

}

func (c *Cpu) CMP(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
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
	data := c.ReadSingleByte(loc)
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
	data := c.ReadSingleByte(loc)
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
	data := c.ReadSingleByte(loc)
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

func (c *Cpu) LSR(mode string, hidden ...*uint8) {
	var data uint8
	var loc uint16
	if mode == "Accumulator" {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.ReadSingleByte(loc)
	}
	if hasBit(data, 0) {
		c.SEC()
	} else {
		c.CLC()
	}
	data = data >> 1
	if mode == "Accumulator" {
		c.aRegister = data
	} else {
		c.WriteSingleByte(loc, data)
	}
	c.ClearNegative()
	if data == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) ASL(mode string, hidden ...*uint8) {
	var data uint8
	var loc uint16
	if mode == "Accumulator" {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.ReadSingleByte(loc)
	}
	if hasBit(data, 7) {
		c.SEC()
	} else {
		c.CLC()
	}
	data = data << 1
	if mode == "Accumulator" {
		c.aRegister = data
	} else {
		c.WriteSingleByte(loc, data)
	}
	c.alterZeroAndNeg(data)
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) ROL(mode string, hidden ...*uint8) {
	var data uint8
	var loc uint16
	if mode == "Accumulator" {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.ReadSingleByte(loc)
	}

	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 7)
	data = data << 1
	if mode == "Accumulator" {
		c.aRegister = data
	} else {
		c.WriteSingleByte(loc, data)
	}
	if temp > 0 {
		data = setBit(data, 0)
	} else {
		data = clearBit(data, 0)
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
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) ROR(mode string, hidden ...*uint8) {
	var data uint8
	var loc uint16
	if mode == "Accumulator" {
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.ReadSingleByte(loc)
	}
	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 0)
	data = data >> 1
	if mode == "Accumulator" {
		c.aRegister = data
	} else {
		c.WriteSingleByte(loc, data)
	}
	if temp > 0 {
		data = setBit(data, 7)
	} else {
		data = clearBit(data, 7)
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
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) JMP(mode string) {

	if mode == "abs" {
		c.pc = c.ReadDoubleByte(c.pc)
	} else {
		loc := c.ReadDoubleByte(c.pc)
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

func (c *Cpu) BMI() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if hasBit(c.statusRegister, 7) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BPL() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if !hasBit(c.statusRegister, 7) {
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVS() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if hasBit(c.statusRegister, 6) {
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVC() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if !hasBit(c.statusRegister, 6) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BCC() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if !hasBit(c.statusRegister, 0) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BEQ() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if hasBit(c.statusRegister, 1) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BCS() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if hasBit(c.statusRegister, 0) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BNE() {
	loc := c.addrMode("imm")
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.ReadSingleByte(loc))
	c.pc++
	if !hasBit(c.statusRegister, 1) {
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) PHA() {
	acc := c.Acc()
	c.Push(acc)
}
func (c *Cpu) PHP() {
	reg := c.statusRegister
	c.Push(reg)
}
func (c *Cpu) PLA() {
	acc := c.Pop()
	if acc == 0 {
		c.SetZero()
	} else {
		if hasBit(acc, 7) {
			c.SetNegative()
		} else {
			c.ClearNegative()
		}
	}
	c.aRegister = acc
}
func (c *Cpu) RTI() {
	c.statusRegister = c.Pop()
	c.pc = c.Pop16()
	if !hasBit(c.statusRegister, 2) {
		c.CLI()
	}

}
func (c *Cpu) BRK() {
	c.Push(c.statusRegister)
	hi := uint8(c.pc >> 8)
	lo := uint8(c.pc & 0b0000000011111111)
	c.Push(hi)
	c.Push(lo)
	c.SEI()
}
func (c *Cpu) JSR() {
	//we need to make sure we increment within the same cycle
	c.Push16(c.pc + 2)
	cal := c.addrMode("imm")
	addr := c.ReadDoubleByte(cal)
	c.pc = addr
}
func (c *Cpu) RTS() {
	val := c.Pop16()
	c.pc = val + 1

}

func (c *Cpu) PLP() {
	reg := c.Pop()
	if reg == 0 {
		c.SetZero()
	} else {
		if hasBit(reg, 7) {
			c.SetNegative()
		} else {
			c.ClearNegative()
		}
	}
	c.statusRegister = reg
}

// func(c *Cpu) DCP(mode string){
// 	loc:=c.addrMode(mode)
// 	data := c.ReadSingleByte(loc)
// 	data--
// 	c.WriteSingleByte(loc,data)

// 	if data == 0 {
// 		c.SetZero()
// 	} else {
// 		c.ClearZero()
// 	}
// 	if hasBit(data, 7) {
// 		c.SetNegative()
// 	} else {
// 		c.ClearNegative()
// 	}

// }
func (c *Cpu) RLA(mode string) {
	var data *uint8
	c.ROL(mode, data)
	c.aRegister = c.aRegister & (*data)
	c.alterZeroAndNeg(c.aRegister)
}
func (c *Cpu) SLA(mode string) {
	var data *uint8
	c.ASL(mode, data)
	c.aRegister = c.aRegister | (*data)
	c.alterZeroAndNeg(c.aRegister)
}
func (c *Cpu) SRE(mode string) {
	var data *uint8
	c.LSR(mode, data)
	c.aRegister = c.aRegister ^ (*data)
	c.alterZeroAndNeg(c.aRegister)

}
func (c *Cpu) AXS(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.xRegister = c.xRegister & c.aRegister
	if data <= c.xRegister {
		c.SEC()
	}
	c.xRegister = c.xRegister - data

	if c.xRegister == 0 {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(c.xRegister, 7) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}
}
func (c *Cpu) ARR(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
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

	if temp == 0 {
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
func (c *Cpu) ANC(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
	c.aRegister = c.aRegister + data
	c.alterZeroAndNeg(c.aRegister)
	if hasBit(c.aRegister, 7) {
		c.SEC()
	} else {
		c.CLC()
	}
}
func (c *Cpu) RRA(mode string) {
	var data *uint8
	c.ROR(mode, data)
	c.aRegister = c.aRegister + *data
	c.alterZeroAndNeg(c.aRegister)
	if hasBit(c.aRegister, 7) {
		c.SEC()
	} else {
		c.CLC()
	}
}
func (c *Cpu) ISB(mode string) {
	var data *uint8
	c.INC(mode, data)
	c.subfromA(*data)

}
func (c *Cpu) subfromA(data uint8) {
	c.aRegister = ((c.aRegister) + uint8(int8(data)-1))
	///lot of questions lol
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
	if hasBit(c.aRegister, 7) {
		c.SEC()
	} else {
		c.CLC()
	}

}

func (c *Cpu) LDA(mode string) {
	loc := c.addrMode(mode)
	data := c.ReadSingleByte(loc)
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

func (c *Cpu) Push(val uint8) {
	loc := STACK + uint16(c.stackPtr)
	c.mem[loc] = val
	c.stackPtr--
}
func (c *Cpu) Push16(val uint16) {
	hi := uint8(val >> 8)
	lo := uint8(val & 0x00FF)
	c.Push(hi)
	c.Push(lo)
}
func (c *Cpu) Pop16() uint16 {
	lo := uint16(c.Pop())
	hi := uint16(c.Pop())
	val := (hi<<8 | lo)
	return val
}
func (c *Cpu) Pop() uint8 {
	c.stackPtr++
	//stack grows down
	loc := STACK + uint16(c.stackPtr)
	temp := c.mem[loc]
	c.mem[loc] = 0
	return temp
}

func (c *Cpu) incrementPassInstruction(inst uint8) {
	inc := uint16(pcIncrement[inst]) - 1
	c.pc = c.pc + inc
}

func (c *Cpu) run() {

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
func (c *Cpu) LoadToMem(data []uint8) {
	copy(c.mem[0xFC:0xFC+len(data)], data)

}

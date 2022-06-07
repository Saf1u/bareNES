package cpu

import (
	"fmt"
)

//should be 0x0600??
const programLocation = 0x0600

const pcStart = 0xFFFC

//Cpu composes of a 6502 register set and addressable memory
type Cpu struct {
	xRegister      uint8
	aRegister      uint8
	yRegister      uint8
	stackPtr       uint8
	pc             uint16
	statusRegister uint8
	cpuBus         bus
}

const STACK uint8 = 0xff

var programLength uint16

const (
	ACCUMULATOR    = "Accumulator"
	IMMEDIATE      = "imm"
	ZERO_PAGE_X    = "zpx"
	ABSOLUTE       = "abs"
	ZERO_PAGE_Y    = "zpy"
	ZERO_PAGE      = "zpg"
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
	INDIRECT       = "ind"
)

const STACK_PAGE uint16 = 0x0100

func (c *Cpu) Acc() uint8 {
	return c.aRegister
}
func (c *Cpu) Stat() uint8 {
	return c.statusRegister
}

func getInst(opcode uint8) string {
	return "LDA"
}
func (c *Cpu) preetyprintSingle() {
	dataLocation := c.pc + 1
	data := c.cpuBus.ReadSingleByte(dataLocation)
	opcode := c.cpuBus.ReadSingleByte(c.pc)
	fmt.Printf("%x	%x %x	%s ", c.pc, opcode, data, getInst(opcode))

}
func (c *Cpu) preetyprintImplied() {
	opcode := c.cpuBus.ReadSingleByte(c.pc)
	fmt.Printf("%x     %x	%s ", c.pc, opcode, getInst(opcode))

}
func (c *Cpu) preetyprintDouble() {
	dataLocation := c.pc + 1
	data := c.cpuBus.ReadSingleByte(dataLocation)
	dataB := c.cpuBus.ReadSingleByte(dataLocation + 1)
	opcode := c.cpuBus.ReadSingleByte(c.pc)
	fmt.Printf("%x	%x %x %x	%s ", c.pc, opcode, data, dataB, getInst(opcode))

}

func (c *Cpu) addrMode(mode string) uint16 {
	var dataLocation uint16
	switch {
	case mode == IMMEDIATE:
		dataLocation = c.pc + 1
		data := c.cpuBus.ReadSingleByte(dataLocation)
		c.preetyprintSingle()
		fmt.Printf("#$%x, ", data)

	case mode == ABSOLUTE:
		c.preetyprintDouble()
		dataLocation = c.cpuBus.ReadDoubleByte(c.pc + 1)
		fmt.Printf("%x,=", dataLocation)
		fmt.Printf("$%x, ", c.cpuBus.ReadSingleByte(dataLocation))

	case mode == ZERO_PAGE_X:
		c.preetyprintSingle()
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.xRegister
		fmt.Printf("$%x,X@ %x,=", data, n)
		dataLocation = uint16(n)
		fmt.Printf("$%x, ", c.cpuBus.ReadSingleByte(dataLocation))

	case mode == ZERO_PAGE_Y:
		c.preetyprintSingle()
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.yRegister
		fmt.Printf("$%x,Y@ %x,=", data, n)
		dataLocation = uint16(n)
		fmt.Printf("$%x, ", c.cpuBus.ReadSingleByte(dataLocation))

	case mode == ABSOLUTE_X:
		c.preetyprintDouble()
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
		fmt.Printf("$%x,X@ %x,=", data, dataLocation)
		fmt.Printf("$%x, ", c.cpuBus.ReadSingleByte(dataLocation))

	case mode == ABSOLUTE_Y:
		c.preetyprintDouble()
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
		fmt.Printf("$%x,Y@ %x,=", data, dataLocation)
		fmt.Printf("$%x, ", c.cpuBus.ReadSingleByte(dataLocation))

	case mode == INDIRECT_X:
		c.preetyprintSingle()
		fmt.Printf("($%x,X) @", c.cpuBus.ReadSingleByte(c.pc+1))
		base := c.cpuBus.ReadSingleByte(c.pc+1) + c.xRegister
		fmt.Printf(" %x =", base)
		low := uint16(c.cpuBus.ReadSingleByte(uint16(base)))
		hi := uint16(c.cpuBus.ReadSingleByte(uint16(base) + 1))
		dataLocation = (hi << 8) | low
		fmt.Printf(" %x =", dataLocation)
		fmt.Printf(" %x", c.cpuBus.ReadSingleByte(dataLocation))
	case mode == INDRECT_Y:
		pos := uint16(c.cpuBus.ReadSingleByte(c.pc + 1))
		low := c.cpuBus.ReadSingleByte(pos)
		hi := c.cpuBus.ReadSingleByte(pos + 1)
		c.preetyprintSingle()
		fmt.Printf("($%x),Y = ", c.cpuBus.ReadSingleByte(c.pc+1))
		loc := uint16(hi)<<8 | uint16(low)
		fmt.Printf("%x @", loc)
		dataLocation = loc + uint16(c.yRegister)
		fmt.Printf("%x = ", dataLocation)
		fmt.Printf("%x", c.cpuBus.ReadSingleByte(dataLocation))
	case mode == ZERO_PAGE:
		c.preetyprintSingle()
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		fmt.Printf("$%x =", data)
		dataLocation = uint16(data)
		fmt.Printf("%x", c.cpuBus.ReadSingleByte(dataLocation))
	case mode == INDIRECT:
		dataLocation = c.cpuBus.ReadDoubleByte(c.pc + 1)

	}
	c.printReg()
	return dataLocation
}

func (c *Cpu) printReg() {
	fmt.Printf("	A:%x X:%x Y:%x P:%x SP:%x", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	fmt.Println()
}

func (c *Cpu) set() {

	c.xRegister = 0
	c.aRegister = 0
	c.yRegister = 0
	c.statusRegister = 0
	c.stackPtr = STACK
	c.cpuBus.WriteDoubleByte(pcStart, programLocation)
	c.pc = c.cpuBus.ReadDoubleByte(pcStart)
}
func (c *Cpu) LoadToMem(data []uint8) {
	c.cpuBus = bus{}
	programLength = programLocation + uint16(len(data))
	copy(c.cpuBus.mem[programLocation:programLength], data)

}
func (c *Cpu) SEC() {
	c.preetyprintImplied()
	c.printReg()
	c.statusRegister = (setBit(c.statusRegister, CARRY_FLAG))

}
func (c *Cpu) CLC() {
	c.preetyprintImplied()
	c.printReg()
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
	c.preetyprintImplied()
	c.printReg()
	c.statusRegister = (setBit(c.statusRegister, INTERRUPT_FLAG))
}
func (c *Cpu) CLI() {
	c.preetyprintImplied()
	c.printReg()
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
	c.preetyprintImplied()
	c.printReg()
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
	data := c.cpuBus.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) LDY(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	c.alterZeroAndNeg(data)

	c.xRegister = data
}
func (c *Cpu) SBC(mode string) {
	//STUBBED
}

func (c *Cpu) ADC(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	temp := c.aRegister + data
	if hasBit(c.statusRegister, CARRY_FLAG) {
		temp++
	}

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
	}
	if hasBit(c.aRegister, 7) {
		c.SetNegative()
	}
}

func (c *Cpu) STA(mode string) {

	loc := c.addrMode(mode)

	c.cpuBus.WriteSingleByte(loc, c.aRegister)
}
func (c *Cpu) STX(mode string) {
	loc := c.addrMode(mode)
	c.cpuBus.WriteSingleByte(loc, c.xRegister)
}
func (c *Cpu) STY(mode string) {
	loc := c.addrMode(mode)
	c.cpuBus.WriteSingleByte(loc, c.yRegister)
}
func (c *Cpu) TAX() {
	c.preetyprintImplied()
	c.printReg()
	data := c.aRegister
	c.alterZeroAndNeg(data)
	c.xRegister = data
}
func (c *Cpu) TAY() {
	c.preetyprintImplied()
	c.printReg()
	data := c.aRegister
	c.alterZeroAndNeg(data)
	c.yRegister = data
}

func (c *Cpu) TXA() {
	c.preetyprintImplied()
	c.printReg()
	data := c.xRegister
	c.alterZeroAndNeg(data)
	c.aRegister = data
}
func (c *Cpu) TYA() {
	c.preetyprintImplied()
	c.printReg()
	data := c.yRegister
	c.alterZeroAndNeg(data)

	c.aRegister = data
}
func (c *Cpu) TXS() {
	c.preetyprintImplied()
	c.printReg()
	data := c.xRegister
	c.stackPtr = data
}

func (c *Cpu) TSX() {
	c.preetyprintImplied()
	c.printReg()
	data := c.stackPtr
	c.alterZeroAndNeg(data)

	c.xRegister = data
}

func (c *Cpu) AND(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	c.aRegister = data & c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) ORA(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	c.aRegister = data | c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) EOR(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	c.aRegister = data ^ c.aRegister
	c.alterZeroAndNeg(c.aRegister)

}

func (c *Cpu) INX() {
	c.preetyprintImplied()
	c.printReg()
	c.xRegister++
	c.alterZeroAndNeg(c.xRegister)

}
func (c *Cpu) INY() {
	c.preetyprintImplied()
	c.printReg()
	c.yRegister++
	c.alterZeroAndNeg(c.yRegister)

}

func (c *Cpu) DEX() {
	c.preetyprintImplied()
	c.printReg()
	c.xRegister--
	c.alterZeroAndNeg(c.xRegister)

}

func (c *Cpu) DEY() {
	c.preetyprintImplied()
	c.printReg()
	c.yRegister--
	c.alterZeroAndNeg(c.yRegister)

}

func (c *Cpu) INC(mode string, hidden ...*uint8) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	data++
	c.cpuBus.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) DEC(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
	data--
	c.cpuBus.WriteSingleByte(loc, data)
	c.alterZeroAndNeg(data)

}

func (c *Cpu) CMP(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
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
	data := c.cpuBus.ReadSingleByte(loc)
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
	data := c.cpuBus.ReadSingleByte(loc)
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
	data := c.cpuBus.ReadSingleByte(loc)
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
	if mode == ACCUMULATOR {
		c.preetyprintImplied()
		fmt.Print("	A")
		c.printReg()
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.cpuBus.ReadSingleByte(loc)
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
		c.cpuBus.WriteSingleByte(loc, data)
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
	if mode == ACCUMULATOR {
		c.preetyprintImplied()
		fmt.Print("	A")
		c.printReg()
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.cpuBus.ReadSingleByte(loc)
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
		c.cpuBus.WriteSingleByte(loc, data)
	}
	c.alterZeroAndNeg(data)
	if len(hidden) != 0 {
		hidden[0] = &data
	}

}

func (c *Cpu) ROL(mode string, hidden ...*uint8) {
	var data uint8
	var loc uint16
	if mode == ACCUMULATOR {
		c.preetyprintImplied()
		fmt.Print("	A")
		c.printReg()
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.cpuBus.ReadSingleByte(loc)
	}

	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 7)
	data = data << 1
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.cpuBus.WriteSingleByte(loc, data)
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
	if mode == ACCUMULATOR {
		c.preetyprintImplied()
		fmt.Print("	A")
		c.printReg()
		data = c.Acc()
	} else {
		loc = c.addrMode(mode)
		data = c.cpuBus.ReadSingleByte(loc)
	}
	temp := c.GetBit(CARRY_FLAG)
	templast := getBit(data, 0)
	data = data >> 1
	if mode == ACCUMULATOR {
		c.aRegister = data
	} else {
		c.cpuBus.WriteSingleByte(loc, data)
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
	if mode == ABSOLUTE {
		c.pc = c.addrMode(ABSOLUTE)
	} else {
		loc := c.cpuBus.ReadDoubleByte(c.pc + 1)
		//6502 HAS A WEIRD WRAPAROUND BUG THAT CAUSES AN ADDRESS TO BE READ BACKWARD IN AN INDIRECT JUMP WE NEED TO REMAIN TRUE TO THIS
		//
		if loc&0x00ff == 0x00ff {
			low := uint16(c.cpuBus.ReadSingleByte(loc))
			hi := uint16(c.cpuBus.ReadSingleByte(loc & 0xFF00))
			fin := hi<<8 | low
			c.pc = fin
		} else {
			c.pc = c.cpuBus.ReadDoubleByte(loc)
		}
	}

}

func (c *Cpu) BMI() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()

	if hasBit(c.statusRegister, 7) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BPL() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()

	if !hasBit(c.statusRegister, 7) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVS() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if hasBit(c.statusRegister, 6) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVC() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if !hasBit(c.statusRegister, 6) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BCC() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if !hasBit(c.statusRegister, 0) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BEQ() {

	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if hasBit(c.statusRegister, 1) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BCS() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if hasBit(c.statusRegister, 0) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BNE() {
	loc := c.pc + 1
	//location of perand to jump too in mem not acc value itself is loc
	toJump := int8(c.cpuBus.ReadSingleByte(loc))
	c.preetyprintSingle()
	fmt.Printf("%x", c.pc+2+uint16(toJump))
	c.printReg()
	if !hasBit(c.statusRegister, 1) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) JSR() {
	//we need to make sure we increment within the same cycle
	c.Push16(c.pc + 3)
	c.preetyprintDouble()
	cal := c.pc + 1
	addr := c.cpuBus.ReadDoubleByte(cal)
	fmt.Printf("%x", addr)
	c.printReg()
	c.pc = addr
}
func (c *Cpu) PHA() {
	c.preetyprintImplied()
	c.printReg()
	acc := c.Acc()
	c.Push(acc)
}
func (c *Cpu) PHP() {
	c.preetyprintImplied()
	c.printReg()
	reg := c.statusRegister
	c.Push(reg)
}
func (c *Cpu) PLA() {
	c.preetyprintImplied()
	c.printReg()
	acc := c.Pop()
	c.alterZeroAndNeg(acc)
	c.aRegister = acc
}
func (c *Cpu) RTI() {
	c.preetyprintImplied()
	c.printReg()
	c.statusRegister = c.Pop()
	c.pc = c.Pop16()
	if !hasBit(c.statusRegister, 2) {
		c.CLI()
	}

}
func (c *Cpu) BRK() {
	c.preetyprintImplied()
	c.printReg()
	c.Push(c.statusRegister)
	c.Push16(c.pc)
	c.SEI()
}

func (c *Cpu) RTS() {
	c.preetyprintImplied()
	c.printReg()
	val := c.Pop16()
	c.pc = val

}

func (c *Cpu) PLP() {
	c.preetyprintImplied()
	reg := c.Pop()
	c.statusRegister = reg
}

/*
Illegal
opcodes to DO


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
*/

func (c *Cpu) LDA(mode string) {
	loc := c.addrMode(mode)
	data := c.cpuBus.ReadSingleByte(loc)
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
	loc := uint16(c.stackPtr)
	c.cpuBus.WriteSingleByte(loc, val)
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
	loc := uint16(c.stackPtr)
	temp := c.cpuBus.ReadSingleByte(loc)
	c.cpuBus.WriteSingleByte(loc, 0)
	return temp
}

type bus struct {
	mem [0xFFFF]uint8
}

//WriteSingleByte writes single byte to mem
func (b *bus) WriteSingleByte(addr uint16, data uint8) {
	//mirror(addr)
	b.mem[addr] = data

}

//ReadSingleByte writes single byte to mem
func (b *bus) ReadSingleByte(addr uint16) uint8 {
	//mirror(addr)
	data := b.mem[addr]
	return data
}

func (b *bus) WriteDoubleByte(addr uint16, data uint16) {
	//mirror(addr)
	low := uint8(data & 0x00FF)
	b.mem[addr] = low
	hi := uint8((data) >> 8)
	b.mem[addr+1] = hi

}
func mirror(addr uint16) uint16 {
	if addr >= 0x000 && addr <= 0x2000 {
		addr = clearBit16(addr, 11)
		addr = clearBit16(addr, 12)
	}
	return addr
}
func (b *bus) ReadDoubleByte(addr uint16) uint16 {
	//mirror(addr)
	var low uint16 = uint16(b.mem[addr])
	var hi uint16 = uint16(b.mem[addr+1])
	res := (hi << 8) | low
	return res
}



func (c *Cpu) Run() {
	c.set()
	for c.pc < programLength {

		temp := c.pc
		location := c.cpuBus.ReadSingleByte(c.pc)
		switch location {
		case 0x00:
			c.BRK()
			return
		case 0x10:
			c.BPL()
		case 0x20:
			c.JSR()
		case 0x30:
			c.BMI()
		case 0x40:
			c.RTI()
		case 0x50:
			c.BVC()
		case 0x60:
			c.RTS()
		case 0x70:
			c.BVC()
		case 0x90:
			c.BCC()
		case 0xA0:
			c.LDY(IMMEDIATE)
		case 0xB0:
			c.BCS()
		case 0xC0:
			c.CPY(IMMEDIATE)
		case 0xD0:
			c.BNE()
		case 0xE0:
			c.CPX(IMMEDIATE)
		case 0xF0:
			c.BEQ()
		case 0x01:
			c.ORA(INDIRECT_X)
		case 0x11:
			c.ORA(INDRECT_Y)
		case 0x21:
			c.AND(INDIRECT_X)
		case 0x31:
			c.AND(INDRECT_Y)
		case 0x41:
			c.EOR(INDIRECT_X)
		case 0x51:
			c.EOR(INDRECT_Y)
		case 0x61:
			c.ADC(INDIRECT_X)
		case 0x71:
			c.ADC(INDRECT_Y)
		case 0x81:
			c.STA(INDIRECT_X)
		case 0x91:
			c.STA(INDRECT_Y)
		case 0xA1:
			c.LDA(INDIRECT_X)
		case 0xB1:
			c.LDA(INDRECT_Y)
		case 0xc1:
			c.CMP(INDIRECT_X)
		case 0xd1:
			c.CMP(INDRECT_Y)
		case 0xe1:
			c.SBC(INDIRECT_X)
		case 0xf1:
			c.SBC(INDRECT_Y)
		case 0xA2:
			c.LDX(IMMEDIATE)
		case 0x24:
			c.BIT(ZERO_PAGE)
		case 0x84:
			c.STY(ZERO_PAGE)
		case 0x94:
			c.STY(ZERO_PAGE_X)
		case 0xA4:
			c.LDY(ZERO_PAGE)
		case 0xB4:
			c.LDY(ZERO_PAGE_X)
		case 0xC4:
			c.CPY(ZERO_PAGE)
		case 0xe4:
			c.CPX(ZERO_PAGE)
		case 0x05:
			c.ORA(ZERO_PAGE)
		case 0x15:
			c.ORA(ZERO_PAGE_X)
		case 0x25:
			c.AND(ZERO_PAGE)
		case 0x35:
			c.AND(ZERO_PAGE_X)
		case 0x45:
			c.EOR(ZERO_PAGE)
		case 0x55:
			c.EOR(ZERO_PAGE_X)
		case 0x65:
			c.ADC(ZERO_PAGE)
		case 0x75:
			c.ADC(ZERO_PAGE_X)
		case 0x85:
			c.STA(ZERO_PAGE)
		case 0x95:
			c.STA(ZERO_PAGE_X)
		case 0xA5:
			c.LDA(ZERO_PAGE)
		case 0xB5:
			c.LDA(ZERO_PAGE_X)
		case 0xc5:
			c.CMP(ZERO_PAGE)
		case 0xd5:
			c.CMP(ZERO_PAGE_X)
		case 0xe5:
			c.SBC(ZERO_PAGE)
		case 0xf5:
			c.SBC(ZERO_PAGE_X)
		case 0x06:
			c.ASL(ZERO_PAGE)
		case 0x16:
			c.ASL(ZERO_PAGE_X)
		case 0x26:
			c.ROL(ZERO_PAGE)
		case 0x36:
			c.ROL(ZERO_PAGE_X)
		case 0x46:
			c.LSR(ZERO_PAGE)
		case 0x56:
			c.LSR(ZERO_PAGE_X)
		case 0x66:
			c.ROR(ZERO_PAGE)
		case 0x76:
			c.ROR(ZERO_PAGE_X)
		case 0x86:
			c.STX(ZERO_PAGE)
		case 0x96:
			c.STX(ZERO_PAGE_Y)
		case 0xA6:
			c.LDX(ZERO_PAGE)
		case 0xB6:
			c.LDX(ZERO_PAGE_Y)
		case 0xc6:
			c.DEC(ZERO_PAGE)
		case 0xd6:
			c.DEC(ZERO_PAGE_X)
		case 0xe6:
			c.INC(ZERO_PAGE)
		case 0xf6:
			c.INC(ZERO_PAGE_X)
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
		case 0x98:
			c.TYA()
		case 0xA8:
			c.TAY()
		case 0xB8:
			c.CLV()
		case 0xc8:
			c.INY()
		case 0xD8:
			//lmao no decimal mode
		case 0xE8:
			c.INX()
		case 0xF8:
			//lmao no decimal mode
		case 0x09:
			c.ORA(IMMEDIATE)
		case 0x19:
			c.ORA(ABSOLUTE_Y)
		case 0x29:
			c.AND(IMMEDIATE)
		case 0x39:
			c.AND(ABSOLUTE_Y)
		case 0x49:
			c.EOR(IMMEDIATE)
		case 0x59:
			c.EOR(ABSOLUTE_Y)
		case 0x69:
			c.ADC(IMMEDIATE)
		case 0x79:
			c.ADC(ABSOLUTE_Y)
		case 0x99:
			c.STA(ABSOLUTE_Y)
		case 0xa9:
			c.LDA(IMMEDIATE)
		case 0xb9:
			c.LDA(ABSOLUTE_Y)
		case 0xc9:
			c.CMP(IMMEDIATE)
		case 0xd9:
			c.CMP(ABSOLUTE_Y)
		case 0xe9:
			c.SBC(IMMEDIATE)
		case 0xf9:
			c.SBC(ABSOLUTE_Y)
		case 0x0a:
			c.ASL(ACCUMULATOR)
		case 0x2a:
			c.ROL(ACCUMULATOR)
		case 0x4a:
			c.LSR(ACCUMULATOR)
		case 0x6a:
			c.ROR(ACCUMULATOR)
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
		case 0xea:
		case 0x2c:
			c.BIT(ABSOLUTE)
		case 0x4c:
			c.JMP(ABSOLUTE)
		case 0x6c:
			c.JMP(INDIRECT)
		case 0x8c:
			c.STY(ABSOLUTE)
		case 0xac:
			c.LDY(ABSOLUTE)
		case 0xbc:
			c.LDY(ABSOLUTE_X)
		case 0xcc:
			c.CPY(ABSOLUTE)
		case 0xec:
			c.CPX(ABSOLUTE)
		case 0x1d:
			c.CPX(ABSOLUTE_X)
		case 0x2d:
			c.AND(ABSOLUTE)
		case 0x3d:
			c.AND(ABSOLUTE_X)
		case 0x4d:
			c.EOR(ABSOLUTE)
		case 0x5d:
			c.EOR(ABSOLUTE_X)
		case 0x6d:
			c.ADC(ABSOLUTE)
		case 0x7d:
			c.ADC(ABSOLUTE_X)
		case 0x8d:
			c.STA(ABSOLUTE)
		case 0x9d:
			c.STA(ABSOLUTE_X)
		case 0xAd:
			c.LDA(ABSOLUTE)
		case 0xBd:
			c.LDA(ABSOLUTE_X)
		case 0xCd:
			c.CMP(ABSOLUTE)
		case 0xDd:
			c.CMP(ABSOLUTE_X)
		case 0xEd:
			c.SBC(ABSOLUTE)
		case 0xFd:
			c.SBC(ABSOLUTE_X)
		case 0x0e:
			c.ASL(ABSOLUTE)
		case 0x1e:
			c.ASL(ABSOLUTE_X)
		case 0x2e:
			c.ROL(ABSOLUTE)
		case 0x3e:
			c.ROL(ABSOLUTE_X)
		case 0x4e:
			c.LSR(ABSOLUTE)
		case 0x5e:
			c.LSR(ABSOLUTE_X)
		case 0x6e:
			c.ROR(ABSOLUTE)
		case 0x7e:
			c.ROR(ABSOLUTE_X)
		case 0x8e:
			c.STX(ABSOLUTE)
		case 0xae:
			c.LDX(ABSOLUTE)
		case 0xbe:
			c.LDX(ABSOLUTE_Y)
		case 0xce:
			c.DEC(ABSOLUTE)
		case 0xde:
			c.DEC(ABSOLUTE_X)
		case 0xee:
			c.INC(ABSOLUTE)
		case 0xfe:
			c.INC(ABSOLUTE_X)
		}

		if c.pc == temp {
			length := pcIncrement[c.cpuBus.ReadSingleByte(c.pc)]
			c.pc = c.pc + (length)

		}
	}

}

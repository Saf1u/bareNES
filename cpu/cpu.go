package cpu

import (
	"fmt"
	"os"
	"strconv"
)

//should be 0x0600??
const programLocation = 0xc000

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

const STACK uint8 = 0xfd

var programLength uint16

const (
	ACCUMULATOR       = "Accumulator"
	RELATIVE          = "relative"
	ABSOLUTE_INDIRECT = "absindirect"
	IMPLIED           = "implied"
	IMMEDIATE         = "imm"
	ZERO_PAGE_X       = "zpx"
	ABSOLUTE          = "abs"
	ZERO_PAGE_Y       = "zpy"
	ZERO_PAGE         = "zpg"
	ABSOLUTE_X        = "absx"
	ABSOLUTE_Y        = "absy"
	INDIRECT_X        = "indx"
	INDIRECT_Y         = "indy"
	CARRY_FLAG        = 0
	ZERO_FLAG         = 1
	INTERRUPT_FLAG    = 2
	DECIMAL_FLAG      = 3
	BREAK_FLAG        = 4
	OVERFLOW_FLAG     = 6
	NEGATIVE_FLAG     = 7
	INDIRECT          = "ind"
)

const STACK_PAGE uint16 = 0x0100

func (c *Cpu) Acc() uint8 {
	return c.aRegister
}
func (c *Cpu) Stat() uint8 {
	return c.statusRegister
}

func getInst(opcode uint8) string {
	return instructionInfo[opcode][0]
}

func (c *Cpu) TraceExecution(mode string) {
	dataLocation := c.pc + 1
	dataSingle := c.cpuBus.ReadSingleByte(dataLocation)
	dataDouble := c.cpuBus.ReadSingleByte(dataLocation + 1)
	opcode := c.cpuBus.ReadSingleByte(c.pc)
	pcDoubleInst := "%04X %02X %02X %02X %s  "
	pcSingleInst := "%04X %02X %02X    %s  "
	pcImpliedInst := "%04X %02X       %s  "
	immediateInstruction := "#$%02X                       "
	absoluteInstruction := "%02X =$%02X                   "
	zeroPageXInstruction := "$%02X,X@ %02X,=$%02X,        "
	zeroPageYInstruction := "$%02X,Y@ %02X,=$%02X,        "
	absoluteXInstruction := "$%02X,X@ %02X,=$%02X,        "
	absoluteYInstruction := "$%02X,Y@ %02X,=$%02X,        "
	indirectInstructionX := "($%02X,X) @ %02X = %02X =%02X"
	indirectInstructionY := "($%02X),Y = %02X @%02X = %02X"
	zeroPageInstruction := "$%02X = %02X                 "
	accumulatorInstruction := "A                            "
	relativeInstruction := "$%04X                        "
	indirectInstruction := "($%04X) = %04X               "

	switch {
	case mode == IMPLIED:
		fmt.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		fmt.Printf("                             A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ACCUMULATOR:
		fmt.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		fmt.Print(accumulatorInstruction)
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == RELATIVE:
		toJump := dataSingle
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		fmt.Printf(relativeInstruction, c.pc+2+uint16(toJump))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == IMMEDIATE:
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		fmt.Printf(immediateInstruction, dataSingle)
		fmt.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ABSOLUTE:
		fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		if (opcode != 0x20) && opcode != 0x4c {
			fmt.Printf(absoluteInstruction, data, c.cpuBus.ReadSingleByte(uint16(data)))
			fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
		} else {
			fmt.Printf(relativeInstruction, data)
			fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
		}

	case mode == ZERO_PAGE_X:
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		var n uint8 = dataSingle + c.xRegister
		dataLocation = uint16(n)
		fmt.Printf(zeroPageXInstruction, dataSingle, n, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ZERO_PAGE_Y:
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		var n uint8 = dataSingle + c.yRegister
		dataLocation = uint16(n)
		fmt.Printf(zeroPageYInstruction, dataSingle, n, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf(":%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ABSOLUTE_X:
		fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
		fmt.Printf(absoluteXInstruction, data, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ABSOLUTE_Y:
		fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
		fmt.Printf(absoluteYInstruction, data, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == INDIRECT_X:
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		base := dataSingle + c.xRegister
		low := uint16(c.cpuBus.ReadSingleByte(uint16(base)))
		hi := uint16(c.cpuBus.ReadSingleByte(uint16(base) + 1))
		dataLocation = (hi << 8) | low
		fmt.Printf(indirectInstructionX, dataSingle, base, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == INDIRECT_Y:
		pos := uint16(dataSingle)
		low := c.cpuBus.ReadSingleByte(pos)
		hi := c.cpuBus.ReadSingleByte(pos + 1)
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)
		fmt.Printf(indirectInstructionY, dataSingle, loc, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == ZERO_PAGE:
		fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		dataLocation = uint16(data)
		fmt.Printf(zeroPageInstruction, dataSingle, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("    A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == INDIRECT:
		fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		loc := c.cpuBus.ReadDoubleByte(c.pc + 1)
		//6502 HAS A WEIRD WRAPAROUND BUG THAT CAUSES AN ADDRESS TO BE READ BACKWARD IN AN INDIRECT JUMP WE NEED TO REMAIN TRUE TO THIS
		//
		if loc&0x00ff == 0x00ff {
			low := uint16(c.cpuBus.ReadSingleByte(loc))
			hi := uint16(c.cpuBus.ReadSingleByte(loc & 0xFF00))
			fin := hi<<8 | low
			fmt.Printf(indirectInstruction, loc, fin)

		} else {
			fin := c.cpuBus.ReadDoubleByte(loc)
			fmt.Printf(indirectInstruction, loc, fin)

		}
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	}

}

func (c *Cpu) addrMode(mode string) uint16 {
	var dataLocation uint16
	switch {
	case mode == RELATIVE:
		loc := c.pc + 1
		toJump := int8(c.cpuBus.ReadSingleByte(loc))
		dataLocation = uint16(toJump)
	case mode == IMMEDIATE:
		dataLocation = c.pc + 1
	case mode == ABSOLUTE:
		dataLocation = c.cpuBus.ReadDoubleByte(c.pc + 1)
	case mode == ZERO_PAGE_X:
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.xRegister
		dataLocation = uint16(n)
	case mode == ZERO_PAGE_Y:
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		var n uint8 = data + c.yRegister
		dataLocation = uint16(n)
	case mode == ABSOLUTE_X:
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
	case mode == ABSOLUTE_Y:
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
	case mode == INDIRECT_X:
		base := c.cpuBus.ReadSingleByte(c.pc+1) + c.xRegister
		low := uint16(c.cpuBus.ReadSingleByte(uint16(base)))
		hi := uint16(c.cpuBus.ReadSingleByte(uint16(base) + 1))
		dataLocation = (hi << 8) | low
	case mode == INDIRECT_Y:
		pos := uint16(c.cpuBus.ReadSingleByte(c.pc + 1))
		low := c.cpuBus.ReadSingleByte(pos)
		hi := c.cpuBus.ReadSingleByte(pos + 1)
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)
	case mode == ZERO_PAGE:
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		dataLocation = uint16(data)
	case mode == INDIRECT:
		dataLocation = c.cpuBus.ReadDoubleByte(c.pc + 1)

	}

	return dataLocation
}

func (c *Cpu) set() {

	c.xRegister = 0
	c.aRegister = 0
	c.yRegister = 0
	c.statusRegister = 0x24
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

	c.yRegister = data
}
func (c *Cpu) SBC(mode string) {
	fmt.Println("sub sus")
	os.Exit(0)

}

func (c *Cpu) ADC(mode string) {
	fmt.Println("add sus")
	os.Exit(0)
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
func (c *Cpu) NOP() {
}

func (c *Cpu) DEY() {
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
			fin := c.cpuBus.ReadDoubleByte(loc)
			c.pc = fin
		}

	}

}

func (c *Cpu) BMI() {
	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 7) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BPL() {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 7) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVS() {
	toJump := c.addrMode(RELATIVE)

	if hasBit(c.statusRegister, 6) {

		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}

func (c *Cpu) BVC() {

	//location of perand to jump too in mem not acc value itself is loc
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 6) {
		c.pc = c.pc + 2
		c.pc = c.pc + (toJump)

	}
}
func (c *Cpu) BCC() {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 0) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BEQ() {

	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 1) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BCS() {
	toJump := c.addrMode(RELATIVE)
	if hasBit(c.statusRegister, 0) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) BNE() {
	toJump := c.addrMode(RELATIVE)
	if !hasBit(c.statusRegister, 1) {
		c.pc = c.pc + 2
		c.pc = c.pc + uint16(toJump)

	}
}
func (c *Cpu) JSR() {
	//we need to make sure we increment within the same cycle
	addr := c.addrMode(ABSOLUTE)
	c.PushDouble(c.pc + 3)
	c.pc = addr
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
	c.alterZeroAndNeg(acc)
	c.aRegister = acc
}
func (c *Cpu) RTI() {
	c.statusRegister = c.Pop()
	c.pc = c.PopDouble()
	if !hasBit(c.statusRegister, 2) {
		c.CLI()
	}

}
func (c *Cpu) BRK() {
	c.Push(c.statusRegister)
	c.PushDouble(c.pc)
	c.SEI()
}

func (c *Cpu) RTS() {
	val := c.PopDouble()
	c.pc = val

}

func (c *Cpu) PLP() {
	reg := c.Pop()

	if hasBit(reg, NEGATIVE_FLAG) {
		c.SetNegative()
	} else {
		c.ClearNegative()
	}

	if hasBit(reg, OVERFLOW_FLAG) {
		c.SetOverflow()
	} else {
		c.statusRegister = clearBit(c.statusRegister, OVERFLOW_FLAG)
	}
	if hasBit(reg, DECIMAL_FLAG) {
		c.SED()
	} else {
		c.CLD()
	}
	if hasBit(reg, INTERRUPT_FLAG) {
		c.SEI()
	} else {
		c.CLI()
	}
	if hasBit(reg, ZERO_FLAG) {
		c.SetZero()
	} else {
		c.ClearZero()
	}
	if hasBit(reg, CARRY_FLAG) {
		c.CLC()
	} else {
		c.SEC()
	}

}

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

//mirror function to map mem locations to similar ranges (to be utilized if i decided to build the nes emulator)
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

func getAddrMode(opcode uint8) string {
	return instructionInfo[opcode][2]
}
func getNumber(opcode uint8) int {
	number := instructionInfo[opcode][1]
	num, err := strconv.Atoi(number)
	if err != nil {
		return -1
	}
	return num
}

func (c *Cpu) Run() {
	c.set()
	for c.pc < programLength {

		temp := c.pc
		location := c.cpuBus.ReadSingleByte(c.pc)
		mode := getAddrMode(location)
		c.TraceExecution(mode)
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
			c.BVS()
		case 0x90:
			c.BCC()
		case 0xA0, 0xA4, 0xB4, 0xAC, 0xBC:
			c.LDY(mode)
		case 0xB0:
			c.BCS()

		case 0xD0:
			c.BNE()

		case 0xF0:
			c.BEQ()
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

		case 0xf5, 0xE9, 0xE5, 0xED, 0xFD, 0xF9, 0xE1, 0xF1:
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
			c.CLD()
			//lmao no decimal mode
		case 0xE8:
			c.INX()
		case 0xF8:
			c.SED()
			//lmao no decimal mode
		case 0x49, 0x45, 0x55, 0x4D, 0x5D, 0x59, 0x41, 0x51:
			c.EOR(mode)

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
		case 0xea:
			c.NOP()
		case 0x2c:
			c.BIT(mode)
		case 0x4c, 0x6C:
			c.JMP(mode)

		case 0x2e:
			c.ROL(mode)
		case 0x3e:
			c.ROL(mode)

		case 0xee, 0xE6, 0xF6, 0xFE:
			c.INC(mode)

		}

		if c.pc == temp {
			val := getNumber(location)
			if val == -1 {
				break
			}
			c.pc = c.pc + uint16(val)

		}
	}

}

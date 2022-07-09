package cpu

import "fmt"

var (
	pcDoubleInst           = "%04X  %02X %02X %02X  %s "
	pcSingleInst           = "%04X  %02X %02X     %s "
	pcImpliedInst          = "%04X  %02X        %s "
	immediateInstruction   = "#$%02X                       "
	absoluteInstruction    = "$%04X = %02X                  "
	zeroPageXInstruction   = "$%02X,X @ %02X = %02X       "
	zeroPageYInstruction   = "$%02X,Y @ %02X = %02X       "
	absoluteXInstruction   = "$%04X,X @ %04X = %02X       "
	absoluteYInstruction   = "$%04X,Y @ %04X = %02X       "
	indirectInstructionX   = "($%02X,X) @ %02X = %04X = %02X"
	indirectInstructionY   = "($%02X),Y = %04X @ %04X = %02X"
	zeroPageInstruction    = "$%02X = %02X                 "
	accumulatorInstruction = "A                           "
	relativeInstruction    = "$%04X                       "
	indirectInstruction    = "($%04X) = %04X              "

	pcDoubleInstIllegal        = "%04X  %02X %02X %02X %s "
	pcSingleInstIllegal        = "%04X  %02X %02X    %s "
	pcImpliedInstIllegal       = "%04X  %02X       %s "
	absoluteInstructionIllegal = "$%04X = %02X                 "
	relativeInstructionIllegal = "$%04X                      "
)

func (c *Cpu) TraceExecution(mode string) {
	dataLocation := c.pc + 1
	dataSingle := c.cpuBus.ReadSingleByte(dataLocation)
	dataDouble := c.cpuBus.ReadSingleByte(dataLocation + 1)
	opcode := c.cpuBus.ReadSingleByte(c.pc)

	switch {
	case mode == IMPLIED:
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcImpliedInstIllegal, c.pc, opcode, getInst(opcode))
		} else {
			fmt.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		}
		fmt.Printf("                            A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ACCUMULATOR:
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcImpliedInstIllegal, c.pc, opcode, getInst(opcode))
		} else {
			fmt.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		}
		fmt.Print(accumulatorInstruction)
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == RELATIVE:
		toJump := dataSingle
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(relativeInstruction, c.pc+2+uint16(int8(toJump)))
		fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
	case mode == IMMEDIATE:
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))

		}
		fmt.Printf(immediateInstruction, dataSingle)
		fmt.Printf(" A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ABSOLUTE:
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
			data := c.cpuBus.ReadDoubleByte(c.pc + 1)
			if (opcode != 0x20) && opcode != 0x4c {
				fmt.Printf(absoluteInstructionIllegal, data, c.cpuBus.ReadSingleByte(uint16(data)))
				fmt.Printf(" A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
			} else {
				fmt.Printf(relativeInstructionIllegal, data)
				fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

			}
		} else {
			fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
			data := c.cpuBus.ReadDoubleByte(c.pc + 1)
			if (opcode != 0x20) && opcode != 0x4c {
				fmt.Printf(absoluteInstruction, data, c.cpuBus.ReadSingleByte(uint16(data)))
				fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
			} else {
				fmt.Printf(relativeInstruction, data)
				fmt.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)
			}
		}

	case mode == ZERO_PAGE_X:
		var n uint8 = dataSingle + c.xRegister
		dataLocation = uint16(n)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(zeroPageXInstruction, dataSingle, n, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("      A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ZERO_PAGE_Y:
		var n uint8 = dataSingle + c.yRegister
		dataLocation = uint16(n)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(zeroPageYInstruction, dataSingle, n, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("      A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ABSOLUTE_X:
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		} else {
			fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		}
		fmt.Printf(absoluteXInstruction, data, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ABSOLUTE_Y:
		data := c.cpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		} else {
			fmt.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))

		}
		fmt.Printf(absoluteYInstruction, data, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == INDIRECT_X:
		base := dataSingle + c.xRegister
		low := uint16(c.cpuBus.ReadSingleByte(uint16(base)))
		temp := uint8(base + 1)
		hi := uint16(c.cpuBus.ReadSingleByte(uint16(temp)))
		dataLocation = (hi << 8) | low
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(indirectInstructionX, dataSingle, base, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("    A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == INDIRECT_Y:
		pos := uint16(dataSingle)
		low := c.cpuBus.ReadSingleByte(pos)
		temp := uint8(dataSingle + 1)
		hi := c.cpuBus.ReadSingleByte(uint16(temp))
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(indirectInstructionY, dataSingle, loc, dataLocation, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ZERO_PAGE:
		data := c.cpuBus.ReadSingleByte(c.pc + 1)
		dataLocation = uint16(data)
		if getInst(opcode)[0:1] == "*" {
			fmt.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			fmt.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		fmt.Printf(zeroPageInstruction, dataSingle, c.cpuBus.ReadSingleByte(dataLocation))
		fmt.Printf("   A:%02X X:%02X Y:%02X P:%02X SP:%02X\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr)

	case mode == ABSOLUTE_INDIRECT:
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

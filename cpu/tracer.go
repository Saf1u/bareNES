package cpu

import (
	"io"
	"log"
	"strconv"
)

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

func (c *Cpu) TraceExecution(mode string, file io.Writer) int {

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(io.Discard)
	tick := 0
	dataLocation := c.pc + 1
	dataSingle := c.CpuBus.ReadSingleByte(dataLocation)
	dataDouble := c.CpuBus.ReadSingleByte(dataLocation + 1)
	opcode := c.CpuBus.ReadSingleByte(c.pc)

	switch {
	case mode == IMPLIED:
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcImpliedInstIllegal, c.pc, opcode, getInst(opcode))
		} else {
			log.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		}
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
		log.Printf("                            A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)

	case mode == ACCUMULATOR:
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcImpliedInstIllegal, c.pc, opcode, getInst(opcode))
		} else {
			log.Printf(pcImpliedInst, c.pc, opcode, getInst(opcode))
		}
		log.Print(accumulatorInstruction)
		log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)

		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == RELATIVE:
		toJump := dataSingle
		if getInst(opcode)[0:1] == "*" {
			//log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			//log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		page := c.pc + 2 + uint16(int8(toJump))
		//log.Printf(relativeInstruction, page)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 3 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		pageEnd := uint16(0x00ff)

		if page > c.pc+2|pageEnd || page < (c.pc+2)&(0xff00) {
			cycleInt += 2
		} else {
			cycleInt++
		}
		tick = cycleInt
	//	log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
	case mode == IMMEDIATE:
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))

		}
		log.Printf(immediateInstruction, dataSingle)
		log.Printf(" A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == ABSOLUTE:
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
			//data := c.CpuBus.ReadDoubleByte(c.pc + 1)
			if (opcode != 0x20) && opcode != 0x4c {
			//	log.Printf(absoluteInstructionIllegal, data, c.CpuBus.ReadSingleByte(uint16(data)))
			//	log.Printf(" A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
			} else {
			//	log.Printf(relativeInstructionIllegal, data)
			//	log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)

			}
		} else {
			log.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
			data := c.CpuBus.ReadDoubleByte(c.pc + 1)
			if (opcode != 0x20) && opcode != 0x4c {
			//	log.Printf(absoluteInstruction, data, c.CpuBus.ReadSingleByte(uint16(data)))
			//	log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
			} else {
				log.Printf(relativeInstruction, data)
			//	log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
			}
		}
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt

	case mode == ZERO_PAGE_X:
		var n uint8 = dataSingle + c.xRegister
		dataLocation = uint16(n)
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		log.Printf(zeroPageXInstruction, dataSingle, n, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("      A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == ZERO_PAGE_Y:
		var n uint8 = dataSingle + c.yRegister
		dataLocation = uint16(n)
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		log.Printf(zeroPageYInstruction, dataSingle, n, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("      A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == ABSOLUTE_X:
		data := c.CpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.xRegister)
		cycle := getCycle(opcode)
		cycleInt := 0

		cycleInt, _ = strconv.Atoi(cycle[0:1])
		if len(cycle) == 2 {
			pageEnd := uint16(0x00ff)

			if dataLocation > data|pageEnd || data > dataLocation {

				cycleInt++
			}
		}
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		} else {
			log.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		}
		log.Printf(absoluteXInstruction, data, dataLocation, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		tick = cycleInt
	case mode == ABSOLUTE_Y:
		data := c.CpuBus.ReadDoubleByte(c.pc + 1)
		dataLocation = data + uint16(c.yRegister)
		cycle := getCycle(opcode)
		cycleInt := 0
		cycleInt, _ = strconv.Atoi(cycle[0:1])
		if len(cycle) == 2 {
			pageEnd := uint16(0x00ff)

			if dataLocation > data|pageEnd || data > dataLocation {

				cycleInt++
			}
		}
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcDoubleInstIllegal, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		} else {
			log.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))

		}
		log.Printf(absoluteYInstruction, data, dataLocation, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		tick = cycleInt
	case mode == INDIRECT_X:
		base := dataSingle + c.xRegister
		low := uint16(c.CpuBus.ReadSingleByte(uint16(base)))
		temp := uint8(base + 1)
		hi := uint16(c.CpuBus.ReadSingleByte(uint16(temp)))
		dataLocation = (hi << 8) | low
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		log.Printf(indirectInstructionX, dataSingle, base, dataLocation, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("    A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == INDIRECT_Y:
		pos := uint16(dataSingle)
		low := c.CpuBus.ReadSingleByte(pos)
		temp := uint8(dataSingle + 1)
		hi := c.CpuBus.ReadSingleByte(uint16(temp))
		loc := uint16(hi)<<8 | uint16(low)
		dataLocation = loc + uint16(c.yRegister)
		cycle := getCycle(opcode)
		cycleInt := 0
		cycleInt, _ = strconv.Atoi(cycle[0:1])
		if len(cycle) == 2 {
			pageEnd := uint16(0x00ff)

			if dataLocation > loc|pageEnd || loc > dataLocation {

				cycleInt++
			}
		}

		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		log.Printf(indirectInstructionY, dataSingle, loc, dataLocation, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("  A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		tick = cycleInt
	case mode == ZERO_PAGE:
		data := c.CpuBus.ReadSingleByte(c.pc + 1)
		dataLocation = uint16(data)
		if getInst(opcode)[0:1] == "*" {
			log.Printf(pcSingleInstIllegal, c.pc, opcode, dataSingle, getInst(opcode))
		} else {
			log.Printf(pcSingleInst, c.pc, opcode, dataSingle, getInst(opcode))
		}
		log.Printf(zeroPageInstruction, dataSingle, c.CpuBus.ReadSingleByte(dataLocation))
		log.Printf("   A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	case mode == ABSOLUTE_INDIRECT:
		log.Printf(pcDoubleInst, c.pc, opcode, dataSingle, dataDouble, getInst(opcode))
		loc := c.CpuBus.ReadDoubleByte(c.pc + 1)
		//6502 HAS A WEIRD WRAPAROUND BUG THAT CAUSES AN ADDRESS TO BE READ BACKWARD IN AN INDIRECT JUMP WE NEED TO REMAIN TRUE TO THIS
		//
		if loc&0x00ff == 0x00ff {
			low := uint16(c.CpuBus.ReadSingleByte(loc))
			hi := uint16(c.CpuBus.ReadSingleByte(loc & 0xFF00))
			fin := hi<<8 | low
			log.Printf(indirectInstruction, loc, fin)

		} else {
			fin := c.CpuBus.ReadDoubleByte(loc)
			log.Printf(indirectInstruction, loc, fin)

		}
		log.Printf("A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d,%3d CYC:%d\n", c.aRegister, c.xRegister, c.yRegister, c.statusRegister, c.stackPtr, c.CpuBus.Ppu.Scanlines, c.CpuBus.Ppu.PpuTicks, c.CpuBus.cpuTicks)
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt
	}
	return tick
}

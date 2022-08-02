package cpu

import "strconv"

var instructionInfo = map[uint8][]string{
	0x69: {"ADC", "2", IMMEDIATE, "2"}, 0x65: {"ADC", "2", ZERO_PAGE, "3"}, 0x75: {"ADC", "2", ZERO_PAGE_X, "4"}, 0x6D: {"ADC", "3", ABSOLUTE, "4"}, 0x7D: {"ADC", "3", ABSOLUTE_X, "4*"}, 0x79: {"ADC", "3", ABSOLUTE_Y, "4*"}, 0x61: {"ADC", "2", INDIRECT_X, "6"}, 0x71: {"ADC", "2", INDIRECT_Y, "5*"},
	0x29: {"AND", "2", IMMEDIATE, "2"}, 0x25: {"AND", "2", ZERO_PAGE, "3"}, 0x35: {"AND", "2", ZERO_PAGE_X, "4"}, 0x2D: {"AND", "3", ABSOLUTE, "4"}, 0x3D: {"AND", "3", ABSOLUTE_X, "4*"}, 0x39: {"AND", "3", ABSOLUTE_Y, "4*"}, 0x21: {"AND", "2", INDIRECT_X, "6"}, 0x31: {"AND", "2", INDIRECT_Y, "5*"},
	0x0A: {"ASL", "1", ACCUMULATOR, "2"}, 0x06: {"ASL", "2", ZERO_PAGE, "5"}, 0x16: {"ASL", "2", ZERO_PAGE_X, "6"}, 0x0E: {"ASL", "3", ABSOLUTE, "6"}, 0x1E: {"ASL", "3", ABSOLUTE_X, "7"},
	0x90: {"BCC", "2", RELATIVE, "2**"},
	0xb0: {"BCS", "2", RELATIVE, "2**"},
	0xF0: {"BEQ", "2", RELATIVE, "2**"},
	0x24: {"BIT", "2", ZERO_PAGE, "3"}, 0x2C: {"BIT", "3", ABSOLUTE, "4"},
	0x30: {"BMI", "2", RELATIVE, "2**"},
	0xD0: {"BNE", "2", RELATIVE, "2**"},
	0x10: {"BPL", "2", RELATIVE, "2**"},
	0x00: {"BRK", "1", IMPLIED, "7"},
	0x50: {"BVC", "2", RELATIVE, "2**"},
	0x70: {"BVS", "2", RELATIVE, "2**"},
	0x18: {"CLC", "1", IMPLIED, "2"},
	0xD8: {"CLD", "1", IMPLIED, "2"},
	0x58: {"CLI", "1", IMPLIED, "2"},
	0xB8: {"CLV", "1", IMPLIED, "2"},
	0xC9: {"CMP", "2", IMMEDIATE, "2"}, 0xC5: {"CMP", "2", ZERO_PAGE, "3"}, 0xD5: {"CMP", "2", ZERO_PAGE_X, "4"}, 0xCD: {"CMP", "3", ABSOLUTE, "4"}, 0xDD: {"CMP", "3", ABSOLUTE_X, "4*"}, 0xD9: {"CMP", "3", ABSOLUTE_Y, "4*"}, 0xC1: {"CMP", "2", INDIRECT_X, "6"}, 0xD1: {"CMP", "2", INDIRECT_Y, "5*"},
	0xE0: {"CPX", "2", IMMEDIATE, "2"}, 0xE4: {"CPX", "2", ZERO_PAGE, "3"}, 0xEC: {"CPX", "3", ABSOLUTE, "4"},
	0xC0: {"CPY", "2", IMMEDIATE, "2"}, 0xC4: {"CPY", "2", ZERO_PAGE, "3"}, 0xCC: {"CPY", "3", ABSOLUTE, "4"},
	0xC6: {"DEC", "2", ZERO_PAGE, "5"}, 0xD6: {"DEC", "2", ZERO_PAGE_X, "6"}, 0xCE: {"DEC", "3", ABSOLUTE, "6"}, 0xDE: {"DEC", "3", ABSOLUTE_X, "7"},
	0xCA: {"DEX", "1", IMPLIED, "2"},
	0x88: {"DEY", "1", IMPLIED, "2"},
	0x49: {"EOR", "2", IMMEDIATE, "2"}, 0x45: {"EOR", "2", ZERO_PAGE, "3"}, 0x55: {"EOR", "2", ZERO_PAGE_X, "4"}, 0x4D: {"EOR", "3", ABSOLUTE, "4"}, 0x5D: {"EOR", "3", ABSOLUTE_X, "4*"}, 0x59: {"EOR", "3", ABSOLUTE_Y, "4*"}, 0x41: {"EOR", "2", INDIRECT_X, "6"}, 0x51: {"EOR", "2", INDIRECT_Y, "5*"},
	0xE6: {"INC", "2", ZERO_PAGE, "5"}, 0xF6: {"INC", "2", ZERO_PAGE_X, "6"}, 0xEE: {"INC", "3", ABSOLUTE, "6"}, 0xFE: {"INC", "3", ABSOLUTE_X, "7"},
	0xE8: {"INX", "1", IMPLIED, "2"},
	0xC8: {"INY", "1", IMPLIED, "2"},
	0x4C: {"JMP", "3", ABSOLUTE, "3"}, 0x6C: {"JMP", "3", ABSOLUTE_INDIRECT, "5"},
	0x20: {"JSR", "3", ABSOLUTE, "6"},
	0xA9: {"LDA", "2", IMMEDIATE, "2"}, 0xA5: {"LDA", "2", ZERO_PAGE, "3"}, 0xB5: {"LDA", "2", ZERO_PAGE_X, "4"}, 0xAD: {"LDA", "3", ABSOLUTE, "4"}, 0xBD: {"LDA", "3", ABSOLUTE_X, "4*"}, 0xB9: {"LDA", "3", ABSOLUTE_Y, "4*"}, 0xA1: {"LDA", "2", INDIRECT_X, "6"}, 0xB1: {"LDA", "2", INDIRECT_Y, "5*"},
	0xA2: {"LDX", "2", IMMEDIATE, "2"}, 0xA6: {"LDX", "2", ZERO_PAGE, "3"}, 0xB6: {"LDX", "2", ZERO_PAGE_Y, "4"}, 0xAE: {"LDX", "3", ABSOLUTE, "4"}, 0xBE: {"LDX", "3", ABSOLUTE_Y, "4*"},
	0xA0: {"LDY", "2", IMMEDIATE, "2"}, 0xA4: {"LDY", "2", ZERO_PAGE, "3"}, 0xB4: {"LDY", "2", ZERO_PAGE_X, "4"}, 0xAC: {"LDY", "3", ABSOLUTE, "4"}, 0xBC: {"LDY", "3", ABSOLUTE_X, "4*"},
	0x4A: {"LSR", "1", ACCUMULATOR, "2"}, 0x46: {"LSR", "2", ZERO_PAGE, "5"}, 0x56: {"LSR", "2", ZERO_PAGE_X, "6"}, 0x4E: {"LSR", "3", ABSOLUTE, "6"}, 0x5E: {"LSR", "3", ABSOLUTE_X, "7"},
	0xEA: {"NOP", "1", IMPLIED, "2"},
	0x09: {"ORA", "2", IMMEDIATE, "2"}, 0x05: {"ORA", "2", ZERO_PAGE, "3"}, 0x15: {"ORA", "2", ZERO_PAGE_X, "4"}, 0x0D: {"ORA", "3", ABSOLUTE, "4"}, 0x1D: {"ORA", "3", ABSOLUTE_X, "4*"}, 0x19: {"ORA", "3", ABSOLUTE_Y, "4*"}, 0x01: {"ORA", "2", INDIRECT_X, "6"}, 0x11: {"ORA", "2", INDIRECT_Y, "5*"},
	0x48: {"PHA", "1", IMPLIED, "3"},
	0x08: {"PHP", "1", IMPLIED, "3"},
	0x68: {"PLA", "1", IMPLIED, "4"},
	0x28: {"PLP", "1", IMPLIED, "4"},
	0x2A: {"ROL", "1", ACCUMULATOR, "2"}, 0x26: {"ROL", "2", ZERO_PAGE, "5"}, 0x36: {"ROL", "2", ZERO_PAGE_X, "6"}, 0x2E: {"ROL", "3", ABSOLUTE, "6"}, 0x3E: {"ROL", "3", ABSOLUTE_X, "7"},
	0x6A: {"ROR", "1", ACCUMULATOR, "2"}, 0x66: {"ROR", "2", ZERO_PAGE, "5"}, 0x76: {"ROR", "2", ZERO_PAGE_X, "6"}, 0x6E: {"ROR", "3", ABSOLUTE, "6"}, 0x7E: {"ROR", "3", ABSOLUTE_X, "7"},
	0x40: {"RTI", "1", IMPLIED, "6"},
	0x60: {"RTS", "1", IMPLIED, "6"},
	0xE9: {"SBC", "2", IMMEDIATE, "2"}, 0xE5: {"SBC", "2", ZERO_PAGE, "3"}, 0xF5: {"SBC", "2", ZERO_PAGE_X, "4"}, 0xED: {"SBC", "3", ABSOLUTE, "4"}, 0xFD: {"SBC", "3", ABSOLUTE_X, "4*"}, 0xF9: {"SBC", "3", ABSOLUTE_Y, "4*"}, 0xE1: {"SBC", "2", INDIRECT_X, "6"}, 0xF1: {"SBC", "2", INDIRECT_Y, "5*"},
	0x38: {"SEC", "1", IMPLIED, "2"},
	0xF8: {"SED", "1", IMPLIED, "2"},
	0x78: {"SEI", "1", IMPLIED, "2"},
	0x85: {"STA", "2", ZERO_PAGE, "3"}, 0x95: {"STA", "2", ZERO_PAGE_X, "4"}, 0x8D: {"STA", "3", ABSOLUTE, "4"}, 0x9D: {"STA", "3", ABSOLUTE_X, "5"}, 0x99: {"STA", "3", ABSOLUTE_Y, "5"}, 0x81: {"STA", "2", INDIRECT_X, "6"}, 0x91: {"STA", "2", INDIRECT_Y, "6"},
	0x86: {"STX", "2", ZERO_PAGE, "3"}, 0x96: {"STX", "2", ZERO_PAGE_Y, "4"}, 0x8E: {"STX", "3", ABSOLUTE, "4"},
	0x84: {"STY", "2", ZERO_PAGE, "3"}, 0x94: {"STY", "2", ZERO_PAGE_X, "4"}, 0x8C: {"STY", "3", ABSOLUTE, "4"},
	0xAA: {"TAX", "1", IMPLIED, "2"},
	0xA8: {"TAY", "1", IMPLIED, "2"},
	0xBA: {"TSX", "1", IMPLIED, "2"},
	0x8A: {"TXA", "1", IMPLIED, "2"},
	0x9A: {"TXS", "1", IMPLIED, "2"},
	0x98: {"TYA", "1", IMPLIED, "2"},
	0x1A: {"*NOP", "1", IMPLIED, "2"},
	0x3A: {"*NOP", "1", IMPLIED, "2"},
	0x5A: {"*NOP", "1", IMPLIED, "2"},
	0x7A: {"*NOP", "1", IMPLIED, "2"},
	0xDA: {"*NOP", "1", IMPLIED, "2"},
	0xFA: {"*NOP", "1", IMPLIED, "2"},
	0xE2: {"*NOP", "2", IMMEDIATE, "2"},
	0x80: {"*NOP", "2", IMMEDIATE, "2"},
	0x82: {"*NOP", "2", IMMEDIATE, "2"},
	0x89: {"*NOP", "2", IMMEDIATE, "2"},
	0xC2: {"*NOP", "2", IMMEDIATE, "2"},
	0x04: {"*NOP", "2", ZERO_PAGE, "3"},
	0x44: {"*NOP", "2", ZERO_PAGE, "3"},
	0x64: {"*NOP", "2", ZERO_PAGE, "3"},
	0x14: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0x34: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0x54: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0x74: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0xD4: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0xF4: {"*NOP", "2", ZERO_PAGE_X, "4"},
	0x0C: {"*NOP", "3", ABSOLUTE, "4"},
	0x1C: {"*NOP", "3", ABSOLUTE_X, "4*"},
	0x3C: {"*NOP", "3", ABSOLUTE_X, "4*"},
	0x5C: {"*NOP", "3", ABSOLUTE_X, "4*"},
	0x7C: {"*NOP", "3", ABSOLUTE_X, "4*"},
	0xDC: {"*NOP", "3", ABSOLUTE_X, "4*"},
	0xFC: {"*NOP", "3", ABSOLUTE_X, "4*"},

	0x02: {"*NOP", "1", IMPLIED, "2"},
	0x12: {"*NOP", "1", IMPLIED, "2"},
	0x22: {"*NOP", "1", IMPLIED, "2"},
	0x32: {"*NOP", "1", IMPLIED, "2"},
	0x42: {"*NOP", "1", IMPLIED, "2"},
	0x52: {"*NOP", "1", IMPLIED, "2"},
	0x62: {"*NOP", "1", IMPLIED, "2"},
	0x72: {"*NOP", "1", IMPLIED, "2"},
	0x92: {"*NOP", "1", IMPLIED, "2"},
	0xB2: {"*NOP", "1", IMPLIED, "2"},
	0xD2: {"*NOP", "1", IMPLIED, "2"},
	0xF2: {"*NOP", "1", IMPLIED, "2"},

	0xA7: {"*LAX", "2", ZERO_PAGE, "3"},
	0xB7: {"*LAX", "2", ZERO_PAGE_Y, "4"},
	0xAF: {"*LAX", "3", ABSOLUTE, "4"},
	0xBF: {"*LAX", "3", ABSOLUTE_Y, "4*"},
	0xA3: {"*LAX", "2", INDIRECT_X, "6"},
	0xB3: {"*LAX", "2", INDIRECT_Y, "5*"},
	0xAB: {"*LXA", "2", IMMEDIATE, "2"},

	0x87: {"*SAX", "2", ZERO_PAGE, "3"},
	0x97: {"*SAX", "2", ZERO_PAGE_Y, "4"},
	0x8f: {"*SAX", "3", ABSOLUTE, "4"},
	0x83: {"*SAX", "2", INDIRECT_X, "6"},
	0xeb: {"*SBC", "2", IMMEDIATE, "2"},

	0xc7: {"*DCP", "2", ZERO_PAGE, "5"},
	0xd7: {"*DCP", "2", ZERO_PAGE_X, "6"},
	0xcf: {"*DCP", "3", ABSOLUTE, "6"},
	0xdf: {"*DCP", "3", ABSOLUTE_X, "7"},
	0xdb: {"*DCP", "3", ABSOLUTE_Y, "7"},
	0xc3: {"*DCP", "2", INDIRECT_X, "8"},
	0xd3: {"*DCP", "2", INDIRECT_Y, "8"},

	0x27: {"*RLA", "2", ZERO_PAGE, "5"},
	0x37: {"*RLA", "2", ZERO_PAGE_X, "6"},
	0x2f: {"*RLA", "3", ABSOLUTE, "6"},
	0x3f: {"*RLA", "3", ABSOLUTE_X, "7"},
	0x3b: {"*RLA", "3", ABSOLUTE_Y, "7"},
	0x23: {"*RLA", "2", INDIRECT_X, "8"},
	0x33: {"*RLA", "2", INDIRECT_Y, "8"},

	0x07: {"*SLO", "2", ZERO_PAGE, "5"},
	0x17: {"*SLO", "2", ZERO_PAGE_X, "6"},
	0x0f: {"*SLO", "3", ABSOLUTE, "6"},
	0x1f: {"*SLO", "3", ABSOLUTE_X, "7"},
	0x1b: {"*SLO", "3", ABSOLUTE_Y, "7"},
	0x03: {"*SLO", "2", INDIRECT_X, "8"},
	0x13: {"*SLO", "2", INDIRECT_Y, "8"},

	0x47: {"*SRE", "2", ZERO_PAGE, "5"},
	0x57: {"*SRE", "2", ZERO_PAGE_X, "6"},
	0x4f: {"*SRE", "3", ABSOLUTE, "6"},
	0x5f: {"*SRE", "3", ABSOLUTE_X, "7"},
	0x5b: {"*SRE", "3", ABSOLUTE_Y, "7"},
	0x43: {"*SRE", "2", INDIRECT_X, "8"},
	0x53: {"*SRE", "2", INDIRECT_Y, "8"},

	0xCB: {"*AXS", "2", IMMEDIATE, "2"},
	0x6B: {"*ARR", "2", IMMEDIATE, "2"},

	0x0B: {"*ANC", "2", IMMEDIATE, "2"},
	0x2B: {"*ANC", "2", IMMEDIATE, "2"},
	0x4B: {"*ALR", "2", IMMEDIATE, "2"},

	0x67: {"*RRA", "2", ZERO_PAGE, "5"},
	0x77: {"*RRA", "2", ZERO_PAGE_X, "6"},
	0x6f: {"*RRA", "3", ABSOLUTE, "6"},
	0x7f: {"*RRA", "3", ABSOLUTE_X, "7"},
	0x7b: {"*RRA", "3", ABSOLUTE_Y, "7"},
	0x63: {"*RRA", "2", INDIRECT_X, "8"},
	0x73: {"*RRA", "2", INDIRECT_Y, "8"},

	0xE7: {"*ISB", "2", ZERO_PAGE, "5"},
	0xF7: {"*ISB", "2", ZERO_PAGE_X, "6"},
	0xEf: {"*ISB", "3", ABSOLUTE, "6"},
	0xFf: {"*ISB", "3", ABSOLUTE_X, "7"},
	0xFb: {"*ISB", "3", ABSOLUTE_Y, "7"},
	0xE3: {"*ISB", "2", INDIRECT_X, "8"},
	0xF3: {"*ISB", "2", INDIRECT_Y, "8"}, //typo from https://www.masswerk.at/6502/6502_instruction_set.html#opcodes-footnote1??

	0x8b: {"*XAA", "2", IMMEDIATE, "2"},
	0xBb: {"*LAS", "3", ABSOLUTE_Y, "4*"},
	0x9b: {"*TAS", "3", ABSOLUTE_Y, "5"},

	0x9f: {"*AHX", "3", ABSOLUTE_Y, "5"},
	0x93: {"*AHX", "2", INDIRECT_Y, "6"},

	0x9E: {"*SHX", "3", ABSOLUTE_Y, "5"},
	0x9c: {"*SHY", "3", ABSOLUTE_X, "5"},
}

func (c *Cpu) getTicks(mode string) int {
	tick := 0
	dataLocation := c.pc + 1
	dataSingle := c.CpuBus.ReadSingleByte(dataLocation)
	opcode := c.CpuBus.ReadSingleByte(c.pc)

	switch mode {

	case ZERO_PAGE, ABSOLUTE_INDIRECT, IMPLIED, ACCUMULATOR, IMMEDIATE, ABSOLUTE, ZERO_PAGE_X, ZERO_PAGE_Y, INDIRECT_X:
		cycle := getCycle(opcode)
		cycleInt := 0
		if len(cycle) == 1 {
			cycleInt, _ = strconv.Atoi(cycle[0:1])
		}
		tick = cycleInt

	case RELATIVE:
		toJump := dataSingle
		page := c.pc + 2 + uint16(int8(toJump))
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

	case ABSOLUTE_X:
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
		tick = cycleInt
	case ABSOLUTE_Y:
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
		tick = cycleInt
	case INDIRECT_Y:
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

		tick = cycleInt
	default:
		panic("unknown mode")

	}
	return tick
}

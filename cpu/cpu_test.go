package cpu

import (
	"testing"

	"github.com/Saf1u/bareNES/utils"
	"github.com/stretchr/testify/assert"
)

func TestReadSingleByte(t *testing.T) {
	var low uint8 = 0b10101011
	cpu := &Cpu{}
	cpu.CpuBus.cpuRam[0] = low
	assert.Equal(t, low, cpu.CpuBus.ReadSingleByte(0))
}

func TestReadDoubleyte(t *testing.T) {
	var low uint8 = 0b10101011
	var hi uint8 = 0b10101010
	cpu := &Cpu{}
	cpu.CpuBus.cpuRam[0] = low
	cpu.CpuBus.cpuRam[1] = hi
	assert.Equal(t, uint16(0b1010101010101011), cpu.CpuBus.ReadDoubleByte(0))
}
func TestWriteSingleByte(t *testing.T) {
	var low uint8 = 0b10101011
	cpu := &Cpu{}
	cpu.CpuBus.WriteSingleByte(0, low)
	assert.Equal(t, low, cpu.CpuBus.cpuRam[0])
}

func TestWriteDoubleyte(t *testing.T) {
	var low uint8 = 0b10101011
	var hi uint8 = 0b10101010
	cpu := &Cpu{}
	cpu.CpuBus.WriteDoubleByte(0, 0b1010101010101011)
	assert.Equal(t, low, cpu.CpuBus.cpuRam[0])
	assert.Equal(t, hi, cpu.CpuBus.cpuRam[1])
}

func TestAddrModes(t *testing.T) {
	t.Run("test relative", func(t *testing.T) {
		val := uint8(0x11)
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.CpuBus.cpuRam[1] = 0x11
		location := cpu.addrMode(RELATIVE)
		assert.Equal(t, location, uint16(val))
	})

	t.Run("test immediate", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		location := cpu.addrMode(IMMEDIATE)
		assert.Equal(t, location, uint16(1))
	})

	t.Run("test absolute", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[2] = 0x31
		location := cpu.addrMode(ABSOLUTE)
		assert.Equal(t, location, uint16(0x3111))
	})
	t.Run("test zpx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		location := cpu.addrMode(ZERO_PAGE_X)
		assert.Equal(t, location, uint16(0x11+0x11))
	})
	t.Run("test zpy", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		location := cpu.addrMode(ZERO_PAGE_Y)
		assert.Equal(t, location, uint16(0x11+0x11))
	})

	t.Run("test absx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[2] = 0x31
		location := cpu.addrMode(ABSOLUTE_X)
		assert.Equal(t, location, uint16(0x3111+0x11))
	})
	t.Run("test absy", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[2] = 0x31
		location := cpu.addrMode(ABSOLUTE_Y)
		assert.Equal(t, location, uint16(0x3111+0x11))
	})

	t.Run("test indx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[0x11+0x11] = 0x31
		cpu.CpuBus.cpuRam[0x11+0x11+1] = 0x21
		location := cpu.addrMode(INDIRECT_X)
		assert.Equal(t, location, uint16(0x2131))
	})

	t.Run("test indY", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[0x11] = 0x31
		cpu.CpuBus.cpuRam[0x11+1] = 0x81
		cpu.CpuBus.cpuRam[0x11+0x11+1] = 0x21
		location := cpu.addrMode(INDIRECT_Y)
		assert.Equal(t, location, uint16(0x8131)+uint16(0x11))
	})

	t.Run("test zeroPage", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.CpuBus.cpuRam[1] = 0x11
		location := cpu.addrMode(ZERO_PAGE)
		assert.Equal(t, location, uint16(0x11))
	})
	t.Run("test indirect", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.CpuBus.cpuRam[1] = 0x11
		cpu.CpuBus.cpuRam[2] = 0x32
		location := cpu.addrMode(INDIRECT)
		assert.Equal(t, location, uint16(0x3211))
	})

}

func TestLDX(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.CpuBus.cpuRam[1] = tt.data
			cpu.LDX(IMMEDIATE)
			assert.Equal(t, cpu.xRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}
func TestLDA(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.CpuBus.cpuRam[1] = tt.data
			cpu.LDA(IMMEDIATE)
			assert.Equal(t, cpu.aRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestLDY(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.CpuBus.cpuRam[1] = tt.data
			cpu.LDY(IMMEDIATE)
			assert.Equal(t, cpu.yRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestStores(t *testing.T) {
	t.Run("sta", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.aRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x05
		cpu.STA(ZERO_PAGE)
		assert.Equal(t, cpu.CpuBus.cpuRam[0x05], uint8(0x11))
	})
	t.Run("stx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x05
		cpu.STX(ZERO_PAGE)
		assert.Equal(t, cpu.CpuBus.cpuRam[0x05], uint8(0x11))
	})
	t.Run("sty", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.CpuBus.cpuRam[1] = 0x05
		cpu.STY(ZERO_PAGE)
		assert.Equal(t, cpu.CpuBus.cpuRam[0x05], uint8(0x11))
	})

}

func TestTAX(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.aRegister = tt.data
			cpu.TAX()
			assert.Equal(t, cpu.xRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestTAY(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.aRegister = tt.data
			cpu.TAY()
			assert.Equal(t, cpu.yRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestTXA(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.xRegister = tt.data
			cpu.TXA()
			assert.Equal(t, cpu.aRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestTYA(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.yRegister = tt.data
			cpu.TYA()
			assert.Equal(t, cpu.aRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestTSX(t *testing.T) {
	var cases = []struct {
		description string
		isZero      bool
		bitSevenOne bool
		data        uint8
	}{
		{"zero", true, false, 0x0},
		{"7th bit negative", false, true, 0b11111111},
		{"7th bit zero", false, false, 0b01111111},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			cpu := &Cpu{}
			cpu.pc = 0
			cpu.stackPtr = tt.data
			cpu.TSX()
			assert.Equal(t, cpu.xRegister, tt.data)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, ZERO_FLAG), tt.isZero)
			assert.Equal(t, utils.HasBit(cpu.statusRegister, NEGATIVE_FLAG), tt.bitSevenOne)
		})
	}
}

func TestTXS(t *testing.T) {
	t.Run("txs", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.xRegister = 0x11
		cpu.TXS()
		assert.Equal(t, cpu.stackPtr, uint8(0x11))
	})

}

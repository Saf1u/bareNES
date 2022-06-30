package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSingleByte(t *testing.T) {
	var low uint8 = 0b10101011
	cpu := &Cpu{}
	cpu.cpuBus.mem[0] = low
	assert.Equal(t, low, cpu.cpuBus.ReadSingleByte(0))
}

func TestReadDoubleyte(t *testing.T) {
	var low uint8 = 0b10101011
	var hi uint8 = 0b10101010
	cpu := &Cpu{}
	cpu.cpuBus.mem[0] = low
	cpu.cpuBus.mem[1] = hi
	assert.Equal(t, uint16(0b1010101010101011), cpu.cpuBus.ReadDoubleByte(0))
}
func TestWriteSingleByte(t *testing.T) {
	var low uint8 = 0b10101011
	cpu := &Cpu{}
	cpu.cpuBus.WriteSingleByte(0, low)
	assert.Equal(t, low, cpu.cpuBus.mem[0])
}

func TestWriteDoubleyte(t *testing.T) {
	var low uint8 = 0b10101011
	var hi uint8 = 0b10101010
	cpu := &Cpu{}
	cpu.cpuBus.WriteDoubleByte(0, 0b1010101010101011)
	assert.Equal(t, low, cpu.cpuBus.mem[0])
	assert.Equal(t, hi, cpu.cpuBus.mem[1])
}

func TestAddrMode(t *testing.T) {
	t.Run("test relative", func(t *testing.T) {
		val := uint8(0x11)
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.cpuBus.mem[1] = 0x11
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
		cpu.cpuBus.mem[1] = 0x11
		cpu.cpuBus.mem[2] = 0x31
		location := cpu.addrMode(ABSOLUTE)
		assert.Equal(t, location, uint16(0x3111))
	})
	t.Run("test zpx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		location := cpu.addrMode(ZERO_PAGE_X)
		assert.Equal(t, location, uint16(0x11+0x11))
	})
	t.Run("test zpy", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		location := cpu.addrMode(ZERO_PAGE_Y)
		assert.Equal(t, location, uint16(0x11+0x11))
	})

	t.Run("test absx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		cpu.cpuBus.mem[2] = 0x31
		location := cpu.addrMode(ABSOLUTE_X)
		assert.Equal(t, location, uint16(0x3111+0x11))
	})
	t.Run("test absy", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		cpu.cpuBus.mem[2] = 0x31
		location := cpu.addrMode(ABSOLUTE_Y)
		assert.Equal(t, location, uint16(0x3111+0x11))
	})

	t.Run("test indx", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.xRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		cpu.cpuBus.mem[0x11+0x11] = 0x31
		cpu.cpuBus.mem[0x11+0x11+1] = 0x21
		location := cpu.addrMode(INDIRECT_X)
		assert.Equal(t, location, uint16(0x2131))
	})

	t.Run("test indY", func(t *testing.T) {
		cpu := &Cpu{}
		cpu.pc = 0
		cpu.yRegister = 0x11
		cpu.cpuBus.mem[1] = 0x11
		cpu.cpuBus.mem[0x11] = 0x31
		cpu.cpuBus.mem[0x11+1] = 0x81
		cpu.cpuBus.mem[0x11+0x11+1] = 0x21
		location := cpu.addrMode(INDIRECT_Y)
		assert.Equal(t, location, uint16(0x8131)+uint16(0x11))
	})
}

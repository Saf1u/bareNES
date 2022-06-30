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


// func TestAddrMode(t *testing.T) {
// 	t.Run("test relative", func(t *testing.T) {
// 		cpu := &Cpu{}
// 		cpu.pc=0
// 		cpu.p
// 	})
// }

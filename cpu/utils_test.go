package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := setBit(num, 1)
	var expected uint8 = 0b10000011
	assert.Equal(t, res, expected)
}

func TestClearBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := clearBit(num, 0)
	var expected uint8 = 0b10000000
	assert.Equal(t, res, expected)
}

func TestHasBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := hasBit(num, 0)
	assert.Equal(t, res, true)
	res = hasBit(num, 1)
	assert.Equal(t, res, false)
}

func TestClearBit16(t *testing.T) {
	var num uint16 = 0b1000000111111111
	res := clearBit16(num, 8)
	var expected uint16 = 0b1000000011111111
	assert.Equal(t, res, expected)
}
func TestGetBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := getBit(num, 0)
	var expected uint8 =0b1
	assert.Equal(t, res, expected)
}

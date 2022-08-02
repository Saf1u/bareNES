package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := SetBit(num, 1)
	var expected uint8 = 0b10000011
	assert.Equal(t, res, expected)
}

func TestClearBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := ClearBit(num, 0)
	var expected uint8 = 0b10000000
	assert.Equal(t, res, expected)
}

func TestHasBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := HasBit(num, 0)
	assert.Equal(t, res, true)
	res = HasBit(num, 1)
	assert.Equal(t, res, false)
}

func TestClearBit16(t *testing.T) {
	var num uint16 = 0b1000000111111111
	res := ClearDoubleByteBit(num, 8)
	var expected uint16 = 0b1000000011111111
	assert.Equal(t, res, expected)
}
func TestGetBit(t *testing.T) {
	var num uint8 = 0b10000001
	res := GetBit(num, 0)
	var expected uint8 = 0b1
	assert.Equal(t, res, expected)
}

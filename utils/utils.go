package utils

func SetBit(num uint8, pos int) uint8 {

	num |= (uint8(1) << pos)
	return num
}

func ClearBit(n uint8, pos int) uint8 {
	var mask uint8 = ^(1 << pos)
	n &= mask
	return n
}
func ClearDoubleByteBit(n uint16, pos int) uint16 {
	var mask uint16 = ^(1 << pos)
	n &= mask
	return n
}

func HasBit(n uint8, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func GetBit(n uint8, pos int) uint8 {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}

func Mirror(addr uint16) uint16 {
	addr = ClearDoubleByteBit(addr, 11)
	addr = ClearDoubleByteBit(addr, 12)
	addr = ClearDoubleByteBit(addr, 13)
	addr = ClearDoubleByteBit(addr, 14)
	addr = ClearDoubleByteBit(addr, 15)
	return addr
}

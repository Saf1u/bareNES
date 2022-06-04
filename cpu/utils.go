package cpu

func setBit(num uint8, pos int) uint8 {

	num |= (uint8(1) << pos)
	return num
}

func clearBit(n uint8, pos int) uint8 {
	var mask uint8 = ^(1 << pos)
	n &= mask
	return n
}
func clearBit16(n uint16, pos int) uint16 {
	var mask uint16 = ^(1 << pos)
	n &= mask
	return n
}

func hasBit(n uint8, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func hasBit16(n uint16, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}
func getBit(n uint8, pos int) uint8 {
	val := n & (1 << pos)
	if val > 0 {
		return 1
	}
	return 0
}

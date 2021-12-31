package cpu

func setBit(num uint8, pos int) uint8 {
	num |= (1 << pos)
	return num
}

func clearBit(n uint8, pos int) uint8 {
	var mask uint8 = ^(1 << pos)
	n &= mask
	return n
}

func hasBit(n uint8, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

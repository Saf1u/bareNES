package main

import (
	"emulator/cpu"
)

func main() {
	// disassembler.ReadFile("code.txt")
	cpu := &cpu.Cpu{}
	// cpu.WriteSingleByte(0x10, 0xFF)
	cpu.LoadToMem([]uint8{0xa5, 0x10, 0x00})
	// fmt.Printf("%b", cpu.Acc())
	// fmt.Println()
	// fmt.Printf("%b", cpu.Stat())
}

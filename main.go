package main

import (
	"emulator/cpu"
)

func main() {
	// disassembler.ReadFile("code.txt")
	cpu := &cpu.Cpu{}
	// cpu.WriteSingleByte(0x10, 0xFF)
	cpu.LoadToMem([]uint8{0x6C, 0xFF, 0x30})
	cpu.Run()
	// fmt.Printf("%b", cpu.Acc())
	// fmt.Println()
	// fmt.Printf("%b", cpu.Stat())

}

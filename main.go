package main

import (
	"emulator/cpu"
)

func main() {
	// disassembler.ReadFile("code.txt")
	cpu := &cpu.Cpu{}
	// cpu.WriteSingleByte(0x10, 0xFF)
	cpu.LoadToMem([]uint8{0xBD, 0x7E, 0x02,0xB5,0x7E,0xB9,0x7E,0x02,0xB9,0x7E,0x00})
	cpu.Run()
	// fmt.Printf("%b", cpu.Acc())
	// fmt.Println()
	// fmt.Printf("%b", cpu.Stat())

}

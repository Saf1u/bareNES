package main

import (
	"emulator/cpu"
	"fmt"
	"io/ioutil"
)

func main() {
	cpu := &cpu.Cpu{}

	content, err := ioutil.ReadFile("nest.nes")
	if err != nil {
		fmt.Println(err)
	}
	//16
	//1662
	max := int(content[4]) * 16384
	rom := []uint8{}
	i := 16

	for ; i < max; i++ {

		rom = append(rom, uint8(content[i]))
	}
	cpu.LoadToMem(rom)
	cpu.Run()
}

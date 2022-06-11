package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	// cpu := &cpu.Cpu{}
	// cpu.LoadToMem([]uint8{0xA9, 0xFF, 0xA2, 0x07, 0x95, 0x10, 0xCA, 0x10, 0xfb})
	// cpu.Run()
	content, err := ioutil.ReadFile("nest.nes")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content[0])
}

package main

import (
	"fmt"
	"emulator/disassembler"
)

func main() {
	res := disassembler.ReadFile("code.txt")
	fmt.Printf("%x",res)
}

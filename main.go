package main

import (
	"bufio"
	"emulator/cpu"
	"emulator/rom"
	"fmt"
	"log"
	"os"
)

func main() {
	cpu := &cpu.Cpu{}
	rom, err := rom.NewRom("nest.nes")
	if err != nil {
		fmt.Println(err)
	}
	cpu.LoadToMem(rom)
	cpu.Run()

	//sys_test()

}

func sys_test() {
	myLog, err := os.Open("mine.log")
	if err != nil {
		log.Println(err)
		return
	}
	rd := bufio.NewReader(myLog)
	myLogConent, err := rd.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	sampleLog, err := os.Open("nestest_no_cycle.log")
	if err != nil {
		log.Println(err)
		return
	}
	rdd := bufio.NewReader(sampleLog)
	sampleLogContent, err := rdd.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	for {
		if myLogConent != sampleLogContent {
			fmt.Println(myLogConent)
			fmt.Println(sampleLogContent)
			os.Exit(0)
		}
		myLogConent, err = rd.ReadString('\n')
		if err != nil {
			break
		}
		sampleLogContent, err = rdd.ReadString('\n')
		if err != nil {
			break
		}
	}

}

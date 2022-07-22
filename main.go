package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// cpu := &cpu.Cpu{}
	// rom, err := rom.NewRom("nest.nes")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// cpu.LoadToMem(rom)
	// cpu.Run()

	sys_test()

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
	sampleLog, err := os.Open("nestest.log")
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
		if myLogConent[0:len(myLogConent)-1] != sampleLogContent[0:len(sampleLogContent)-2] {
			myLogConent = myLogConent[0 : len(myLogConent)-1]
			fmt.Println(myLogConent)
			sampleLogContent = sampleLogContent[0 : len(sampleLogContent)-2]
			fmt.Println(sampleLogContent)
			// for i := 0; i < len(myLogConent); i++ {
			// 	if myLogConent[i] != sampleLogContent[i] {
			// 		fmt.Println(string(myLogConent[i]))
			// 		fmt.Println(string(sampleLogContent[i-1]))
			// 		fmt.Println(i)
			// 	}
			// }
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

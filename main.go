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

	rom, err := rom.NewRom("nest.nes")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	//game := rom.CharRom
	cpu := &cpu.Cpu{}

	cpu.LoadToMem(rom)

	cpu.Run()
}

// func sdlI(frame []uint8) {
// 	err := sdl.InitSubSystem(sdl.INIT_VIDEO)
// 	defer sdl.Quit()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	window, err := sdl.CreateWindow("sdl window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256*3, 240*3, sdl.WINDOW_SHOWN)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	err = renderer.SetScale(3, 3)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STATIC, 256, 240)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	err = texture.Update(nil, frame, 256*3)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = renderer.Copy(texture, nil, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(0)
// 	}
// 	renderer.Present()
// 	sdl.PumpEvents()
// 	for true {
// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch event.(type) {
// 			case *sdl.QuitEvent:
// 				return

// 			}

// 		}
// 	}
// }

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

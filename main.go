package main

import (
	"emulator/cpu"
	"emulator/rom"
	"fmt"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	rom, err := rom.NewRom("nest.nes")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	renderEvent := make(chan int)
	cpu := &cpu.Cpu{}
	cpu.Init(rom, renderEvent)
	go func() {
		cpu.Run()
	}()

	Screen(cpu.CpuBus.Ppu.Frame.Screen[:], renderEvent)
}

func Screen(pixels []uint8, renderLisitiner chan int) {
	err := sdl.InitSubSystem(sdl.INIT_VIDEO)
	defer sdl.Quit()
	if err != nil {
		log.Println(err)
		return
	}
	window, err := sdl.CreateWindow("NES", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256*3, 240*3, sdl.WINDOW_SHOWN)
	defer window.Destroy()
	if err != nil {
		log.Println(err)
		return
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC)
	defer renderer.Destroy()
	if err != nil {
		log.Println(err)
		return
	}

	err = renderer.SetScale(3, 3)
	if err != nil {
		log.Println(err)
		return
	}
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STATIC, 256, 240)
	defer texture.Destroy()
	if err != nil {
		fmt.Println(err)
	}

	sdl.PumpEvents()
	for {
		<-renderLisitiner
		err = texture.Update(nil, pixels, 256*3)
		if err != nil {
			fmt.Println(err)
		}
		err = renderer.Copy(texture, nil, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return

			}

		}

	}

}

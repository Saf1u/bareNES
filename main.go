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

	rom, err := rom.NewRom("roms/cyo.nes")

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	renderEvent := make(chan int)
	done := make(chan int)
	cpu := &cpu.Cpu{}
	cpu.Init(rom, renderEvent, done)
	go func() {
		cpu.Run()
	}()

	Screen(cpu.CpuBus.Ppu.Frame.Screen[:], renderEvent, done, &cpu.CpuBus.Pad)
}

func Screen(pixels []uint8, renderLisitiner chan int, done chan int, pad *cpu.JoyPad) {
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
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:

				switch t.Keysym.Sym {
				case sdl.K_UP:
					if t.State == sdl.PRESSED {
						//fmt.Println("up pressed")
						pad.Set(true, 4)
					}
					if t.State == sdl.RELEASED {
						//fmt.Println("up released")
						pad.Set(false, 4)
					}
				case sdl.K_DOWN:
					if t.State == sdl.PRESSED {
						//fmt.Println("down pressed")
						pad.Set(true, 5)
					}
					if t.State == sdl.RELEASED {
						//fmt.Println("down released")
						pad.Set(false, 5)
					}
				case sdl.K_RIGHT:
					if t.State == sdl.PRESSED {
						//fmt.Println("down pressed")
						pad.Set(true, 7)
					}
					if t.State == sdl.RELEASED {
						//fmt.Println("down released")
						pad.Set(false, 7)
					}
				case sdl.K_LEFT:
					if t.State == sdl.PRESSED {
						//fmt.Println("down pressed")
						pad.Set(true, 6)
					}
					if t.State == sdl.RELEASED {
						//fmt.Println("down released")
						pad.Set(false, 6)
					}
				case sdl.K_a:
					if t.State == sdl.PRESSED {
						//fmt.Println("down pressed")
						pad.Set(true, 0)
					}
					if t.State == sdl.RELEASED {
						//fmt.Println("down released")
						pad.Set(false, 0)
					}
				case sdl.K_b:
					if t.State == sdl.PRESSED {
						pad.Set(true, 1)
					}
					if t.State == sdl.RELEASED {
						pad.Set(false, 1)
					}
				case sdl.K_RETURN:
					if t.State == sdl.PRESSED {
						pad.Set(true, 3)
					}
					if t.State == sdl.RELEASED {
						pad.Set(false, 3)
					}
				case sdl.K_SPACE:
					if t.State == sdl.PRESSED {
						pad.Set(true, 2)
					}
					if t.State == sdl.RELEASED {
						pad.Set(false, 2)
					}
				}

			}

		}
		done <- 0
	}

}

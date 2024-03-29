
# BareNES

## Build instructions
I havent tested on windows yet, but the following steps should be sufficent to run on mac devices
### Depndencies
use brew to install go and sdl:
```bash
brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
brew install go
```
#### Option A:
clone the repo and cd into directory and build:
```bash
git clone https://github.com/Saf1u/bareNES
cd bareNES
go build .
```
Only 2 roms are currently supported, nestest.rom (a component test for cpu) and cyo (a homebrew game more details below on it)
```bash
./bareNES roms/nestest.nes 
./bareNES roms/cyo.nes 
```
#### Option B:
Simply download the 2 rom files, ensure you have installed sdl and go, then run 
``` bash
go get github.com/Saf1u/bareNES@latest
```
which will install the binary in your gopath ~/go/bin directory as BareNes. You can then:
```bash
bareNES [downloaded-rom-file]
```
to run if you have your gopath in your system path variable or

```bash
~/go/bin/bareNES [downloaded-rom-file]
```
if your path is not in the env variable

button mappings

| conroller     | Keyboard      |    
| ------------- |:-------------:| 
| Action button    | A|
| Up      | Up key      |
| Down| Down Key      | 
| left      | left key      |
| right| right Key      | 
| select| spacebar     | 
| start| enter key     | 


## Description
A nintendo entertaiment system emulator written in go. The goal of this project was to gain an understanding of how old computer systems without operating systems and all the fancy things we have today worked to achieve their functionality. During the course of developing the emulator, I gained a better insight into the following technologies and mechanisms:

1. CPU instruction interpretation and execution (Particularly 6502 assembly but the skills I learned here are transferable I believe)
2. Importance of Interrupts and clock cycle management between "distributed" componenets in a computer system
3. Hexadecimal/binary arithmetic and general bit manipulation algoritihms
4. Memory adressing modes
5. Subroutine calls/conditional branches and system stack management 
6. 8bit graphics manimpulation (Understanding how color and location are encoded in memory to render the pixelated characters on our screen)
7. 2D image rendering using the sdl library 
8. And much more....

Its amazing how much devs of the past were able to accomplish with such meager tools.

## Helpful resources if you like this stuff too
The following are a list of resources that I found useful in implementing the emulator:
1. [Middle Engine's 6502 blog](https://www.middle-engine.com/blog/posts/2020/06/23/programming-the-nes-the-6502-in-detail) I recommend using this to study the opcode behaviours, adressing modes etc to familiarize yourself with the cpu.
2. [Masswerk 6502 instruction set](https://www.masswerk.at/6502/6502_instruction_set.html) may be a little intimidating if you are new to this stuff but it's an excellent resource that goes in depth on the opcodes and instruction behaviour in general.
3. [Nes Dev wiki](https://www.nesdev.org/wiki/Nesdev_Wiki) an mazing resource that goes into detail about the system. Although extremely helpful, i found it a bit difficult siftting through hardware specifities and information that a game dev would need to know versus an emulator dev, nevertheless still a very crucial resource especially when it comes to understanding ppu (picture processing uinit) semantics.
4. [Bugzmanov's Ebook on implementing a nes emulator in rust](https://bugzmanov.github.io/nes_ebook/chapter_1.html) Highly recommend this to supplement the sometimes cryptic info on nesdev. He has excellent explanantions of emulator behaviour plus code examples to clarify ambiguities.

## Current stage
The emulator is not 100% complete, but you could run the homebrew game [cyo](https://www.nesworld.com/article.php?system=nes&data=neshomebrew) which is included in the repo for a demo.I'm still fixing some errors as I go, and making optimizations in the implementation but I am just greatful to be at this point lol. The CPU passes the famed nestest which I also included as a rom file that can be run similarly to a game (i'd say its a good test of cpu behaviour and controller/mem mappings but not necessarily ppu behaviour)

## TODO'S:
- [ ] Implement scrolling in ppu
- [ ] write more unit tests to improve robustness
- [ ] Implement APU (Audio processing unit)
- [ ] Repay tech debt (Probably will be a continous process)
- [ ] Make the emulator available on homebrew for mac users
- [ ] test on windows.

 



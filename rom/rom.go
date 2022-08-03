package rom

import (
	"emulator/common"
	"errors"
	"io/ioutil"
)

var (
	nesTags      = [4]uint8{0x4E, 0x45, 0x53, 0x1A}
	progRom uint = 16384
	charRom uint = 8192
)

type Rom struct {
	ProgramRom []uint8
	CharRom    []uint8
	Mapper     uint8
	MirorType  int
}

func NewRom(file string) (*Rom, error) {
	rom := &Rom{}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 4; i++ {
		if content[i] != nesTags[i] {
			return nil, errors.New("not in ines")
		}
	}
	ctrlByteOne := content[6]
	mapperLow := (ctrlByteOne >> 4)
	ctrlByteTwo := content[7]
	mapperHigh := (ctrlByteTwo) & 0b11110000
	mapper := (mapperHigh) | mapperLow
	rom.Mapper = mapper
	inesVer := ctrlByteTwo >> 2 & 0b11111111
	if inesVer != 0 {
		return nil, errors.New("not 1.0")
	}
	if ctrlByteOne&0b010000000 != 0 {
		rom.MirorType = common.FOUR_SCREEN
	} else {
		if ctrlByteOne&0b1 == 0 {
			rom.MirorType = common.HORIZONTAL
		} else {
			rom.MirorType = common.VERTICAL
		}
	}
	maxProg := uint(content[4]) * progRom

	maxChar := uint(content[5]) * charRom

	i := uint(16)
	if (content[6])&0b100 == 0 {
		i = uint(16 + 512)
	}
	i = uint(16)

	for ; i < maxProg+16; i++ {

		rom.ProgramRom = append(rom.ProgramRom, uint8(content[i]))
	}
	for i := maxProg + 16; i < maxProg+maxChar; i++ {

		rom.CharRom = append(rom.CharRom, uint8(content[i]))
	}
	return rom, nil

}

func (r *Rom) ReadRom(addr uint16) uint8 {
	addr = addr - 0x8000
	if len(r.ProgramRom) <= 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	data := r.ProgramRom[addr]
	return data

}

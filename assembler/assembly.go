package assembler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseToOpcodes(inst string, hexaDump *mem) {

	switch {
	//lda immediate
	case regex("(?i)^LDA\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}

		hexaDump.mem[instLoader] = 0xa9
		instLoader++

		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break
		//lda zeropage
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xa5
		instLoader++

		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break
		//lda zeropage x
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xb5
		instLoader++

		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break
		//lda abs
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xad
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xbd
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xb9
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^LDA\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xa1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^LDA\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xb1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^ADC\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x69
		instLoader++

		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^ADC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x65
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^ADC\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x75
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^ADC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x6d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^ADC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x7d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^ADC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x79
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^ADC\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x71
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^ADC\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x61
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^ASL\\s*A$", inst):
		hexaDump.mem[instLoader] = 0x0A
		instLoader++
		break

		//lda zeropage
	case regex("(?i)^ASL\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x06
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^ASL\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x16
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^ASL\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x0E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^ASL\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x1e
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^AND\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}

		hexaDump.mem[instLoader] = 0x1e
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^AND\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x25
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^AND\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x35
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^AND\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x2d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^AND\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x3d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^AND\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x39
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^AND\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x31
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^AND\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x21
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^BCS\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xB0
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^BCC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x90
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^BEQ\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xF0
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^BIT\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x24
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^BIT\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x2c
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break

	case regex("(?i)^BMI\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x30
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^BPL\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x10
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^BRK\\s*$", inst):

		hexaDump.mem[instLoader] = 0x00
		instLoader++
		break

	case regex("(?i)^BVS\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x70
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^BVC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x50
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^CLC\\s*$", inst):
		hexaDump.mem[instLoader] = 0x18
		instLoader++
		break
	case regex("(?i)^CLD\\s*$", inst):
		hexaDump.mem[instLoader] = 0xd8
		instLoader++
		break
	case regex("(?i)^CLI\\s*$", inst):
		hexaDump.mem[instLoader] = 0x58
		instLoader++
		break
	case regex("(?i)^CLV\\s*$", inst):
		hexaDump.mem[instLoader] = 0xB8
		instLoader++
		break

	case regex("(?i)^CMP\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC9
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^CMP\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC5
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^CMP\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xD5
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^CMP\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xCd
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^CMP\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xDd
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^CMP\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xD9
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^CMP\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xD1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^CMP\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^CPX\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE0
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^CPX\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE4
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs
	case regex("(?i)^CPX\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xEC
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^CPY\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC0
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^CPY\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC4
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs
	case regex("(?i)^CPY\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xCC
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^DEC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xC6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^DEC\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xD6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^DEC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xCE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^DEC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xDE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^DEX\\s*$", inst):

		hexaDump.mem[instLoader] = 0xCA
		instLoader++
		break

	case regex("(?i)^DEY\\s*$", inst):

		hexaDump.mem[instLoader] = 0x88
		instLoader++
		break

	case regex("(?i)^EOR\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x49
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^EOR\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x45
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^EOR\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x55
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^LDA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x4d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^EOR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x5d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^EOR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x59
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^EOR\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x51
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^EOR\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x41
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^INC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^INC\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xF6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^INC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xFE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^INC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xFE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^INX\\s*$", inst):

		hexaDump.mem[instLoader] = 0xE8
		instLoader++
		break

	case regex("(?i)^INY\\s*$", inst):

		hexaDump.mem[instLoader] = 0xC8
		instLoader++
		break

	case regex("(?i)^JMP\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x4C
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^JMP\\s*\\(\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x6C
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
	case regex("(?i)^JSR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x20
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda immediate
	case regex("(?i)^LDX\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xa2
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^LDX\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xa6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^LDX\\s*\\$[0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xb6
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^LDX\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xaE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^LDX\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xbE
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^LDY\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xa0
		instLoader++

		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^LDY\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xA4
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^LDY\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xb4
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^LDY\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xaC
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^LDY\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*X$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xbC
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^LSR\\sA*$", inst):
		fmt.Println("LSR", inst)
		hexaDump.mem[instLoader] = 0x4A
		instLoader++
		break

	case regex("(?i)^LSR\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x46
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^LSR\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x56
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^LSR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x4E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^LSR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*X$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x5E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^NOP\\s*$", inst):

		hexaDump.mem[instLoader] = 0xEA
		instLoader++
		break

	case regex("(?i)^ORA\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x09
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage
	case regex("(?i)^ORA\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x05
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^ORA\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x15
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^ORA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x0d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^ORA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x1d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^ORA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x19
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^ORA\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x11
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^ORA\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x01
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^PHA\\s*$", inst):

		hexaDump.mem[instLoader] = 0x48
		instLoader++
		break
	case regex("(?i)^PHP\\s*$", inst):

		hexaDump.mem[instLoader] = 0x08
		instLoader++
		break
	case regex("(?i)^PLA\\s*$", inst):

		hexaDump.mem[instLoader] = 0x68
		instLoader++
		break
	case regex("(?i)^PLP\\s*$", inst):

		hexaDump.mem[instLoader] = 0x28
		instLoader++
		break

	case regex("(?i)^ROL\\sA*$", inst):

		hexaDump.mem[instLoader] = 0x2A
		instLoader++
		break

	case regex("(?i)^ROL\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x26
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^ROL\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x36
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^ROL\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x2E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^ROL\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*X$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		hexaDump.mem[instLoader] = 0x3E
		instLoader++

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^ROR\\sA*$", inst):

		hexaDump.mem[instLoader] = 0x6A
		instLoader++
		break

	case regex("(?i)^ROR\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x66
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^ROR\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x76
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^ROR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x6E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^ROR\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*X$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x7E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^RTI\\s*$", inst):

		hexaDump.mem[instLoader] = 0x40
		instLoader++
		break
	case regex("(?i)^RTS\\s*$", inst):

		hexaDump.mem[instLoader] = 0x60
		instLoader++
		break

	case regex("(?i)^SBC\\s*#\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)#\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseUint(val[2:], 16, 8)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE9
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break
		//lda zeropage
	case regex("(?i)^SBC\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE5
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break
		//lda zeropage x
	case regex("(?i)^SBC\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xF5
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^SBC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xEd
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^SBC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xFd
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^SBC\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		hexaDump.mem[instLoader] = 0xF9
		instLoader++
		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^SBC\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xF1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//ind x
	case regex("(?i)^SBC\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0xE1
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^STA\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x85
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^STA\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x95
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^STA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x8d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs x
	case regex("(?i)^STA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x9d
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

		//lda abs y
	case regex("(?i)^STA\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x99
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//indirect y
	case regex("(?i)^STA\\s*\\(\\$[0-9a-f][0-9a-f]\\),\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x91
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		break

		//ind x
	case regex("(?i)^STA\\s*\\(\\$[0-9a-f][0-9a-f],\\s*x\\)$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x81
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^SEC\\s*$", inst):

		hexaDump.mem[instLoader] = 0x38
		instLoader++
		break

	case regex("(?i)^SED\\s*$", inst):

		hexaDump.mem[instLoader] = 0xF8
		instLoader++
		break

	case regex("(?i)^SEI\\s*$", inst):

		hexaDump.mem[instLoader] = 0x78
		instLoader++
		break

	case regex("(?i)^STX\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		hexaDump.mem[instLoader] = 0x86
		instLoader++
		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^STX\\s*\\$[0-9a-f][0-9a-f],\\s*y$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)
		hexaDump.mem[instLoader] = 0x96
		instLoader++
		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^STX\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x8E
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^STY\\s*\\$[0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x84
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda zeropage x
	case regex("(?i)^STY\\s*\\$[0-9a-f][0-9a-f],\\s*x$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x94
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break
		//lda abs
	case regex("(?i)^STY\\s*\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]$", inst):
		re := regexp.MustCompile("(?i)\\$[0-9a-f][0-9a-f][0-9a-f][0-9a-f]")
		str := re.Find([]byte(inst))
		val := string(str)

		i, err := strconv.ParseInt(val[1:3], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = 0x8C
		instLoader++
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++
		i, err = strconv.ParseInt(val[3:], 16, 16)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		hexaDump.mem[instLoader] = uint8(i)
		instLoader++

		break

	case regex("(?i)^TAX\\s*$", inst):

		hexaDump.mem[instLoader] = 0xAA
		instLoader++
		break

	case regex("(?i)^TAY\\s*$", inst):

		hexaDump.mem[instLoader] = 0xA8
		instLoader++
		break
	case regex("(?i)^TSX\\s*$", inst):

		hexaDump.mem[instLoader] = 0xBA
		instLoader++
		break
	case regex("(?i)^TXA\\s*$", inst):

		hexaDump.mem[instLoader] = 0x8A
		instLoader++
		break
	case regex("(?i)^TXS\\s*$", inst):

		hexaDump.mem[instLoader] = 0x9A
		instLoader++
		break
	case regex("(?i)^TYA\\s*$", inst):

		hexaDump.mem[instLoader] = 0x98
		instLoader++
		break
	}

}

type mem struct {
	mem [1000]uint8
}

var instLoader = 0

func regex(pattern string, inst string) bool {
	res, _ := regexp.MatchString(pattern, inst)
	return res
}

var instructions = []string{}

//ReadFile reads instructions provided as memonics and returns hex mapping
func ReadFile(filename string) []uint8 {
	var hexaDump = &mem{}
	fl, err := os.OpenFile(filename, os.O_RDWR, 755)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	defer fl.Close()
	reader := bufio.NewReader(fl)
	parseProgram(reader)
	for _, c := range instructions {
		parseToOpcodes(c, hexaDump)
	}

	return hexaDump.mem[0:instLoader]
}

func parseProgram(reader *bufio.Reader) {

	for {
		rd, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if rd == "\n" {
			continue
		}
		match, _ := regexp.MatchString("^;.", rd)
		if match {
			continue
		}

		processInstruction(rd)

	}

}

func processInstruction(token string) {
	instruction := strings.Split(token, ";")
	instruction[0] = strings.TrimSpace(instruction[0])
	instructions = append(instructions, instruction[0])

}

package cpu_test

import (
	"emulator/cpu"
	"testing"

	"github.com/stretchr/testify/suite"
)

type cpuTestSuite struct {
	suite.Suite
	cpu *cpu.Cpu
}

func (suite *cpuTestSuite) SetupTest() {
	suite.cpu = &cpu.Cpu{}
	suite.cpu.LoadToMem([]uint8{0x10, 0x20, 0x30, 0xFF, 0x20})
}
func (suite *cpuTestSuite) TestGetBit() {
	suite.Equal(uint8(0), suite.cpu.GetBit(1))

}

func (suite *cpuTestSuite) TestZeroFlag() {
	suite.cpu.SetZero()
	suite.Equal(uint8(1), suite.cpu.GetBit(1))
}

func (suite *cpuTestSuite) TestLDA() {
suite.cpu.LDA(cpu.IMMEDIATE)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(cpuTestSuite))
}

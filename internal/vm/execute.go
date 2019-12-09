package vm

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/onlyafly/oakblue/internal/spec"
)

const (
	pc_start   = 0x3000 // default PC start location
	mask_11111 = 0x1F   // = 11111
	mask_111   = 0x7    // = 0111
	//TODO mask_11    = 0x3    // = 0011
	mask_1 = 0x1 // = 0001

	memory_size = math.MaxUint16
)

type Machine struct {
	// The VM has 65,536 memory locations, each of which stores a 16-bit value
	mem [memory_size]uint16

	regs [spec.MaxRegisters]uint16
}

func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) RegisterDump() string {
	var b strings.Builder

	for i, reg := range m.regs {
		if i == 8 {
			b.WriteString(fmt.Sprintf("PC=%#v ", reg))
		} else if i == 9 {
			b.WriteString(fmt.Sprintf("COND=%#v", reg))
		} else {
			b.WriteString(fmt.Sprintf("R%d=%#v ", i, reg))
		}
	}

	return b.String()
}

func (m *Machine) LoadMemory(data []byte, loadAddress uint16) {
	im := loadAddress
	for id := 0; id+1 < len(data); id += 2 {
		m.mem[im] = binary.BigEndian.Uint16(data[id : id+2])
		im++
	}
}

func (m *Machine) Execute() {

	// set the PC to starting position
	m.regs[spec.R_PC] = pc_start

	running := true
	for running {

		// ORDERING: The PC must only be incremented after its use is complete
		if m.regs[spec.R_PC] >= memory_size {
			return // end of memory reached
		}
		instr := m.readMemory(m.regs[spec.R_PC])
		m.regs[spec.R_PC]++

		op := instr >> 12

		switch op {
		case spec.OP_ADD:
			// ADD
			//  15-12  opcode
			//  11-09  DR: destination register
			//  08-06  SR1: sum register 1
			//  05     mode: 0 = register, 1 = immediate
			//  04-03  (if mode=0) 00
			//  02-00  (if mode=0) SR2: sum register 2
			//  04-00  (if mode=1) IMM5: immediate value, sign extended

			dr := (instr >> 9) & mask_111
			sr1 := (instr >> 6) & mask_111
			mode := (instr >> 5) & mask_1

			if mode == 1 {
				imm5 := signExtend(instr&mask_11111, 5)
				m.regs[dr] = m.regs[sr1] + imm5
			} else {
				sr2 := instr & mask_111
				m.regs[dr] = m.regs[sr1] + m.regs[sr2]
			}

			m.updateFlags(dr)
		case spec.OP_AND:
			// FIXME
		case spec.OP_NOT:
			// FIXME
		case spec.OP_BR:
			// FIXME
		case spec.OP_JMP:
			// FIXME
		case spec.OP_JSR:
			// FIXME
		case spec.OP_LD:
			// FIXME
		case spec.OP_LDI:
			// FIXME
		case spec.OP_LDR:
			// FIXME
		case spec.OP_LEA:
			// FIXME
		case spec.OP_ST:
			// FIXME
		case spec.OP_STI:
			// FIXME
		case spec.OP_STR:
			// FIXME
		case spec.OP_TRAP:
			// FIXME
		case spec.OP_RES:
			// FIXME
		case spec.OP_RTI:
			// FIXME
		default:
			// FIXME
		}
	}

	fmt.Println("Hello, Oakblue")
}

func (m *Machine) readMemory(loc uint16) uint16 {
	return m.mem[loc]
}

// TODO write a unit test for this
// Any time a value is written to a register, we need to update the flags to indicate its sign
func (m *Machine) updateFlags(r uint16) {
	if m.regs[r] == 0 {
		m.regs[spec.R_COND] = spec.FL_ZRO
	} else if (m.regs[r] >> 15) == 1 { // a 1 in the left-most bit indicates negative
		m.regs[spec.R_COND] = spec.FL_NEG
	} else {
		m.regs[spec.R_COND] = spec.FL_POS
	}
}

// TODO write a unit test for this
func signExtend(x uint16, bitCount int) uint16 {
	if ((x >> (bitCount - 1)) & 1) == 1 {
		x |= (0xFFFF << bitCount)
	}
	return x
}

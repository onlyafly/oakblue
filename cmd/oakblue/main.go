package main

import "fmt"
import "math"

// The VM has 65,536 memory locations, each of which stores a 16-bit value
var memory [math.MaxUint16]uint16

// Registers:
//  8 general purpose registers (R0-R7)
//  1 program counter (PC) register
//  1 condition flags (COND) register
const (
	r_r0 = iota
	r_r1
	r_r2
	r_r3
	r_r4
	r_r5
	r_r6
	r_r7
	r_pc
	r_cond
	maxRegisters
)

// Each instruction is 16 bits long, with the left 4 bits storing the opcode. The rest of the bits
// are used to store the parameters.
const (
	op_br   = iota // branch
	op_add         // add
	op_ld          // load
	op_st          // store
	op_jsr         // jump register
	op_and         // bitwise and
	op_ldr         // load register
	op_str         // store register
	op_rti         // unused
	op_not         // bitwise not
	op_ldi         // load indirect
	op_sti         // store indirect
	op_jmp         // jump
	op_res         // reserved (unused)
	op_lea         // load effective address
	op_trap        // execute trap
)

var regs [maxRegisters]uint16

const (
	fl_pos = 1 << iota // P
	fl_zro             // Z
	fl_neg             // N
)

const (
	pc_start   = 0x3000 // default PC start location
	mask_11111 = 0x1F   // = 11111
	mask_111   = 0x7    // = 0111
	mask_11    = 0x3    // = 0011
	mask_1     = 0x1    // = 0001
)

func main() {

	// set the PC to starting position
	regs[r_pc] = pc_start

	running := true
	for running {
		regs[r_pc]++

		instr := readMemory(regs[r_pc])
		op := instr >> 12

		switch op {
		case op_add:
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
				regs[dr] = regs[sr1] + imm5
			} else {
				sr2 := instr & mask_111
				regs[dr] = regs[sr1] + regs[sr2]
			}

			updateFlags(dr)
		case op_and:
			// FIXME
		case op_not:
			// FIXME
		case op_br:
			// FIXME
		case op_jmp:
			// FIXME
		case op_jsr:
			// FIXME
		case op_ld:
			// FIXME
		case op_ldi:
			// FIXME
		case op_ldr:
			// FIXME
		case op_lea:
			// FIXME
		case op_st:
			// FIXME
		case op_sti:
			// FIXME
		case op_str:
			// FIXME
		case op_trap:
			// FIXME
		case op_res:
			// FIXME
		case op_rti:
			// FIXME
		default:
			// FIXME
		}
	}

	fmt.Println("Hello, Oakblue")
}

func readMemory(loc uint16) uint16 {
	return memory[loc]
}

// TODO write a unit test for this
func signExtend(x uint16, bitCount int) uint16 {
	if ((x >> (bitCount - 1)) & 1) == 1 {
		x |= (0xFFFF << bitCount)
	}
	return x
}

// TODO write a unit test for this
// Any time a value is written to a register, we need to update the flags to indicate its sign
func updateFlags(r uint16) {
	if regs[r] == 0 {
		regs[r_cond] = fl_zro
	} else if (regs[r] >> 15) == 1 { // a 1 in the left-most bit indicates negative
		regs[r_cond] = fl_neg
	} else {
		regs[r_cond] = fl_pos
	}
}

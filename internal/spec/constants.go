package spec

// Registers:
//  8 general purpose registers (R0-R7)
//  1 program counter (PC) register
//  1 condition flags (COND) register
const (
	R_R0 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC
	R_COND
	MaxRegisters
)

var RegisterNames = [...]string{
	"R0",
	"R1",
	"R2",
	"R3",
	"R4",
	"R5",
	"R6",
	"R7",
	"PC",
	"COND",
}

// Each instruction is 16 bits long, with the left 4 bits storing the opcode. The rest of the bits
// are used to store the parameters.
const (
	OP_BR   = iota // branch
	OP_ADD         // add
	OP_LD          // load
	OP_ST          // store
	OP_JSR         // jump register
	OP_AND         // bitwise and
	OP_LDR         // load register
	OP_STR         // store register
	OP_RTI         // unused
	OP_NOT         // bitwise not
	OP_LDI         // load indirect
	OP_STI         // store indirect
	OP_JMP         // jump
	OP_RES         // reserved (unused)
	OP_LEA         // load effective address
	OP_TRAP        // execute trap
)

var OpcodeNames = [...]string{
	"BR",
	"ADD",
	"LD",
	"ST",
	"JSR",
	"AND",
	"LDR",
	"STR",
	"RTI",
	"NOT",
	"LDI",
	"STI",
	"JMP",
	"RES",
	"LEA",
	"TRAP",
}

const (
	FL_POS = 1 << iota // P
	FL_ZRO             // Z
	FL_NEG             // N
)

const (
	TRAPVECT_GETC  = 0x20
	TRAPVECT_OUT   = 0x21
	TRAPVECT_PUTS  = 0x22
	TRAPVECT_IN    = 0x23
	TRAPVECT_PUTSP = 0x24
	TRAPVECT_HALT  = 0x25
)

const (
	DefaultOrigin = 0x3000
	PCStart       = 0x3000 // default PC start location
)

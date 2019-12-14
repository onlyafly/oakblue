package emitter

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/stretchr/testify/assert"
)

func TestEmit_Add(t *testing.T) {
	program := ast.NewProgram([]ast.Statement{
		&ast.Instruction{
			Opcode: spec.OP_ADD,
			Dr:     spec.R_R7,
			Sr1:    spec.R_R2,
			Mode:   1,
			Imm5:   1,
		},
	}, ast.NewSymbolTable(), 0x3000)

	actual, err := Emit(program, syntax.NewErrorList("Emit"))
	assert.NoError(t, err)

	expected := []byte{0x30, 0x0, 0x1e, 0xa1}
	assert.EqualValues(t, expected, actual)

	/*
		instr := binary.BigEndian.Uint16(actual)

		mask_11111 := uint16(0x1F) // = 11111
		mask_111 := uint16(0x7)    // = 0111
		//mask_11 := uint16(0x3)     // = 0011
		mask_1 := uint16(0x1) // = 0001
		dr := (instr >> 9) & mask_111
		sr1 := (instr >> 6) & mask_111
		mode := (instr >> 5) & mask_1
		imm5 := instr & mask_11111

		assert.Equal(t, "", fmt.Sprintf("%d %d %d %d", dr, sr1, mode, imm5))
	*/
}

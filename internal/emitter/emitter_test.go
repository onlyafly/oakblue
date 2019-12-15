package emitter

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/stretchr/testify/assert"
)

func TestEmit_Add1(t *testing.T) {
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
}

func TestEmit_Add2(t *testing.T) {
	program := ast.NewProgram([]ast.Statement{
		&ast.Instruction{
			Opcode: spec.OP_ADD,
			Dr:     spec.R_R0,
			Sr1:    spec.R_R0,
			Mode:   1,
			Imm5:   16,
		},
	}, ast.NewSymbolTable(), 0x3000)

	actual, err := Emit(program, syntax.NewErrorList("Emit"))
	assert.NoError(t, err)

	expected := []byte{
		0x30, 0x0, // Header
		0b00010000, 0b00110000,
	}
	assert.EqualValues(t, expected, actual)
}

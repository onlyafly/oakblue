package ast

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/stretchr/testify/assert"
)

func TestProgram_String_OneStatement(t *testing.T) {
	p := NewProgram([]Statement{
		&Instruction{
			Opcode: spec.OP_ADD,
			Dr:     spec.R_R0,
			Sr1:    spec.R_R0,
			Mode:   1,
			Imm5:   1,
		},
	})

	expected := "ADD R0 R0 1"
	actual := p.String()

	assert.Equal(t, expected, actual)
}

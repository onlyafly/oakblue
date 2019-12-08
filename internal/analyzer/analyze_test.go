package analyzer

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/cst"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {

	input := cst.Listing([]*cst.Line{
		cst.NewLine([]cst.Node{
			cst.NewSymbol("ADD"),
			cst.NewSymbol("R0"),
			cst.NewSymbol("R0"),
			cst.NewInteger(1),
		}),
	})

	actual, err := Analyze(input, syntax.NewErrorList("Syntax"))

	if assert.NoError(t, err) {
		expected := ast.NewProgram([]ast.Statement{
			&ast.Instruction{
				Opcode: spec.OP_ADD,
				Dr:     spec.R_R0,
				Sr1:    spec.R_R0,
				Mode:   1,
				Imm5:   1,
			},
		})

		assert.EqualValues(t, expected, actual)
	}

}

func TestAnalyze_Add_NotEnoughArgs(t *testing.T) {

	input := cst.Listing([]*cst.Line{
		cst.NewLine([]cst.Node{
			cst.NewSymbol("ADD"),
			cst.NewSymbol("R0"),
			cst.NewInteger(1),
		}),
	})

	_, err := Analyze(input, syntax.NewErrorList("Syntax"))
	assert.Error(t, err)
}

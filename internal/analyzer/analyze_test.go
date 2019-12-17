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
			cst.NewSymbol(".ORIG"),
			cst.NewHexNumber(0x3000),
		}),
		cst.NewLine([]cst.Node{
			cst.NewSymbol("ADD"),
			cst.NewRegister(spec.R_R0),
			cst.NewRegister(spec.R_R0),
			cst.NewDecimalNumber(1),
		}),
	})

	actual, err := Analyze(input, syntax.NewErrorList("Syntax"))
	if !assert.NoError(t, err) {
		return
	}

	expected := ast.NewProgram([]ast.Statement{
		&ast.Instruction{
			Opcode: spec.OP_ADD,
			Dr:     spec.R_R0,
			Sr1:    spec.R_R0,
			Mode:   1,
			Imm5:   1,
		},
	}, ast.NewSymbolTable(), 0x3000)

	assert.EqualValues(t, expected, actual)
}

func TestAnalyze_Add(t *testing.T) {

	input := cst.Listing([]*cst.Line{
		cst.NewLine([]cst.Node{
			cst.NewSymbol("ADD"),
			cst.NewRegister(spec.R_R0),
			cst.NewRegister(spec.R_R0),
			cst.NewHexNumber(0xf),
		}),
	})

	actual, err := Analyze(input, syntax.NewErrorList("Syntax"))
	if !assert.NoError(t, err) {
		return
	}

	expected := ast.NewProgram([]ast.Statement{
		&ast.Instruction{
			Opcode: spec.OP_ADD,
			Dr:     spec.R_R0,
			Sr1:    spec.R_R0,
			Mode:   1,
			Imm5:   15,
		},
	}, ast.NewSymbolTable(), 0x0)

	assert.EqualValues(t, expected, actual)
}

func TestAnalyze_Add_NotEnoughArgs(t *testing.T) {

	input := cst.Listing([]*cst.Line{
		cst.NewLine([]cst.Node{
			cst.NewSymbol("ADD"),
			cst.NewRegister(spec.R_R0),
			cst.NewDecimalNumber(1),
		}),
	})

	_, err := Analyze(input, syntax.NewErrorList("Syntax"))
	assert.Error(t, err)
}

func Test_analyzer_analyzeRegister(t *testing.T) {
	a := &analyzer{errors: syntax.NewErrorList("analysis")}

	actual := a.analyzeRegister(cst.NewRegister(spec.R_R0))
	expected := spec.R_R0
	assert.Equal(t, expected, actual)

	actual = a.analyzeRegister(cst.NewRegister(spec.R_R7))
	expected = spec.R_R7
	assert.Equal(t, expected, actual)

	a.analyzeRegister(cst.NewRegister(spec.R_COND))
	assert.Error(t, a.errors)
}

func Test_analyzer_analyzeNumber(t *testing.T) {

	a := &analyzer{errors: syntax.NewErrorList("analysis")}

	want := 0
	n := cst.NewDecimalNumber(16)
	got := a.analyzeNumber(n, "TEST1", 5)
	assert.Equal(t, want, got)

	want = 0
	n = cst.NewDecimalNumber(-17)
	got = a.analyzeNumber(n, "TEST1", 5)
	assert.Equal(t, want, got)

	want = 15
	n = cst.NewDecimalNumber(15)
	got = a.analyzeNumber(n, "TEST1", 5)
	assert.Equal(t, want, got)

	want = -16
	n = cst.NewDecimalNumber(-16)
	got = a.analyzeNumber(n, "TEST1", 5)
	assert.Equal(t, want, got)

	want = 7
	h := cst.NewHexNumber(0b00000111)
	got = a.analyzeNumber(h, "TEST1", 3)
	assert.Equal(t, want, got)

	want = 0
	h = cst.NewHexNumber(0b00001111)
	got = a.analyzeNumber(h, "TEST1", 3)
	assert.Equal(t, want, got)
}

package parser

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/stretchr/testify/assert"
)

func TestParseSymbol(t *testing.T) {
	errors := syntax.NewErrorList("Syntax")

	result := parseSymbol(Token{Value: "fred"}, errors)
	assert.Equal(t, "fred", result.String())
}

func TestParse_Simple(t *testing.T) {
	input := `ADD R0 R0 1`
	result, err := Parse(input, "test", syntax.NewErrorList("Syntax"))
	if assert.NoError(t, err) {
		assert.Equal(t, "ADD R0 R0 1", result.String())
	}
}

func TestParse_TwoLines(t *testing.T) {
	input := `
	ADD R0 R0 1
	ADD R1 R1 1
	`
	result, err := Parse(input, "test", syntax.NewErrorList("Syntax"))
	if assert.NoError(t, err) {
		assert.Equal(t, "ADD R0 R0 1\nADD R1 R1 1", result.String())
	}
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSymbol(t *testing.T) {
	errors := NewParserErrorList()

	result := parseSymbol(Token{Value: "fred"}, errors)
	assert.Equal(t, "fred", result.String())
}

func TestParse_Simple(t *testing.T) {
	input := `ADD R0 R0 1`
	result, err := Parse(input, "test")
	if assert.NoError(t, err) {
		assert.Equal(t, "ADD R0 R0 1", result.String())
	}
}

func TestParse_TwoStatements(t *testing.T) {
	input := `
	ADD R0 R0 1
	ADD R1 R1 1
	`
	result, err := Parse(input, "test")
	if assert.NoError(t, err) {
		assert.Equal(t, "ADD R0 R0 1\nADD R1 R1 1", result.String())
	}
}

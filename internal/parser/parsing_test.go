package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSymbol(t *testing.T) {
	errors := NewParserErrorList()

	result := parseSymbol(Token{Value: "fred"}, &errors)
	assert.Equal(t, "fred", result.String())
}

func TestParse_Simple(t *testing.T) {
	result, _ := Parse("ADD R0 R0 1", "test")
	assert.Equal(t, "ADD R0 R0 1", result.String())
}

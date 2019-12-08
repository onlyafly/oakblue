package cst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListing_String_OneStatement(t *testing.T) {
	p := Listing([]*Statement{
		NewStatement([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
	})

	expected := "ADD R0 R0 1"
	actual := p.String()

	assert.Equal(t, expected, actual)
}

func TestListing_String_MultipleStatements(t *testing.T) {
	p := Listing([]*Statement{
		NewStatement([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
		NewStatement([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
		NewStatement([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
	})

	expected := "ADD R0 R0 1\nADD R0 R0 1\nADD R0 R0 1"
	actual := p.String()

	assert.Equal(t, expected, actual)
}

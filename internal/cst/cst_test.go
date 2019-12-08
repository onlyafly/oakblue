package cst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListing_String_OneLine(t *testing.T) {
	p := Listing([]*Line{
		NewLine([]Node{
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

func TestListing_String_MultipleLines(t *testing.T) {
	p := Listing([]*Line{
		NewLine([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
		NewLine([]Node{
			NewSymbol("ADD"),
			NewSymbol("R0"),
			NewSymbol("R0"),
			NewInteger(1),
		}),
		NewLine([]Node{
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

package parser

import (
	"github.com/onlyafly/oakblue/internal/util"
	"testing"
)

func TestParseSymbol(t *testing.T) {
	errors := NewParserErrorList()

	result1 := parseSymbol(Token{Value: "fred"}, &errors)
	util.CheckEqualString(t, "fred", result1.String())
}

func TestParse_Simple(t *testing.T) {
	result, _ := Parse("ADD R0 R0 1", "test")

	util.CheckEqualString(t, "ADD R0 R0 1", result.String())
}

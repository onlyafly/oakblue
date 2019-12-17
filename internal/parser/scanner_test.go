package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScannerToString(t *testing.T) {
	scanner, tokens := Scan("testing", "ADD R2 R0 R1")
	assert.Equal(t, "ADD", (<-tokens).String())

	actual := scanner.String()
	expected := "<scanner next=\"R2\">"
	assert.Equal(t, expected, actual)
}

func TestScan_Symbols(t *testing.T) {
	_, tokens := Scan("testing", "ADD R2 R0 R1 ROUGH")

	tok := <-tokens
	assert.Equal(t, "ADD", tok.String())
	assert.Equal(t, TcSymbol, tok.Code)

	tok = <-tokens
	assert.Equal(t, "R2", tok.String())
	assert.Equal(t, TcRegister, tok.Code)

	tok = <-tokens
	assert.Equal(t, "R0", tok.String())
	assert.Equal(t, TcRegister, tok.Code)

	tok = <-tokens
	assert.Equal(t, "R1", tok.String())
	assert.Equal(t, TcRegister, tok.Code)

	tok = <-tokens
	assert.Equal(t, "ROUGH", tok.String())
	assert.Equal(t, TcSymbol, tok.Code)
}

func TestScan_Newline(t *testing.T) {
	_, tokens := Scan("testing", "1\n2")
	assert.Equal(t, "1", (<-tokens).String())
	assert.Equal(t, "NEWLINE", (<-tokens).String())
	assert.Equal(t, "2", (<-tokens).String())
}

func TestScan_Numbers(t *testing.T) {
	_, tokens := Scan("testing", "5 -5 #5 #-5 xf0 0xf0")

	tok := <-tokens
	assert.Equal(t, "5", tok.String())
	assert.Equal(t, TcDecimalNumber, tok.Code)

	tok = <-tokens
	assert.Equal(t, "-5", tok.String())
	assert.Equal(t, TcDecimalNumber, tok.Code)

	tok = <-tokens
	assert.Equal(t, "#5", tok.String())
	assert.Equal(t, TcDecimalNumber, tok.Code)

	tok = <-tokens
	assert.Equal(t, "#-5", tok.String())
	assert.Equal(t, TcDecimalNumber, tok.Code)

	tok = <-tokens
	assert.Equal(t, "xf0", tok.String())
	assert.Equal(t, TcHexNumber, tok.Code)

	tok = <-tokens
	assert.Equal(t, "0xf0", tok.String())
	assert.Equal(t, TcHexNumber, tok.Code)
}

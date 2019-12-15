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
	_, tokens := Scan("testing", "ADD R2 R0 R1")

	assert.Equal(t, "ADD", (<-tokens).String())
	assert.Equal(t, "R2", (<-tokens).String())
	assert.Equal(t, "R0", (<-tokens).String())
	assert.Equal(t, "R1", (<-tokens).String())
}

func TestScan_Numbers(t *testing.T) {
	_, tokens := Scan("testing", "1")
	assert.Equal(t, "1", (<-tokens).String())
}

func TestScan_Newline(t *testing.T) {
	_, tokens := Scan("testing", "1\n2")
	assert.Equal(t, "1", (<-tokens).String())
	assert.Equal(t, "NEWLINE", (<-tokens).String())
	assert.Equal(t, "2", (<-tokens).String())
}

func TestScan_Hex(t *testing.T) {
	_, tokens := Scan("testing", "x1")

	tok := <-tokens
	assert.Equal(t, "x1", tok.String())
	assert.Equal(t, TcHexNumber, tok.Code)
}

/* TODO: delete me

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	assert.Equal(t, "(", (<-tokens).String())
	assert.Equal(t, "1", (<-tokens).String())
	assert.Equal(t, "2", (<-tokens).String())
	assert.Equal(t, "3", (<-tokens).String())
	assert.Equal(t, ")", (<-tokens).String())
	assert.Equal(t, "EOF", (<-tokens).String())

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	assert.Equal(t, "(", (<-tokens).String())
	assert.Equal(t, "abc", (<-tokens).String())
	assert.Equal(t, "ab2?", (<-tokens).String())
	assert.Equal(t, "3.5", (<-tokens).String())
	assert.Equal(t, ")", (<-tokens).String())
	assert.Equal(t, "EOF", (<-tokens).String())

	_, tokens = Scan("tester3", "\\a")
	assert.Equal(t, "\\a", (<-tokens).String())

	fmt.Printf("START\n")
	_, tokens = Scan("tester4", "(list \\a)")
	assert.Equal(t, "(", (<-tokens).String())
	assert.Equal(t, "list", (<-tokens).String())
	assert.Equal(t, "\\a", (<-tokens).String())
	assert.Equal(t, ")", (<-tokens).String())
	fmt.Printf("END\n")
}
*/

package parser

import (
	"testing"

	"github.com/onlyafly/oakblue/internal/util"
)

func TestScannerToString(t *testing.T) {
	scanner, tokens := Scan("testing", "ADD R2 R0 R1")
	util.CheckEqualStringer(t, "ADD", <-tokens)

	actual := scanner.String()
	expected := "<scanner next=\"R2\">"
	util.CheckEqualStringer(t, expected, actual)
}

func TestScan_Symbols(t *testing.T) {
	_, tokens := Scan("testing", "ADD R2 R0 R1")

	util.CheckEqualStringer(t, "ADD", <-tokens)
	util.CheckEqualStringer(t, "R2", <-tokens)
	util.CheckEqualStringer(t, "R0", <-tokens)
	util.CheckEqualStringer(t, "R1", <-tokens)
}

func TestScan_Numbers(t *testing.T) {
	_, tokens := Scan("testing", "1")
	util.CheckEqualStringer(t, "1", <-tokens)
}

func TestScan_Newline(t *testing.T) {
	_, tokens := Scan("testing", "1\n2")
	util.CheckEqualStringer(t, "1", <-tokens)
	util.CheckEqualStringer(t, "NEWLINE", <-tokens)
	util.CheckEqualStringer(t, "2", <-tokens)
}

/* TODO delete me

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	util.CheckEqualStringer(t, "(", <-tokens)
	util.CheckEqualStringer(t, "1", <-tokens)
	util.CheckEqualStringer(t, "2", <-tokens)
	util.CheckEqualStringer(t, "3", <-tokens)
	util.CheckEqualStringer(t, ")", <-tokens)
	util.CheckEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	util.CheckEqualStringer(t, "(", <-tokens)
	util.CheckEqualStringer(t, "abc", <-tokens)
	util.CheckEqualStringer(t, "ab2?", <-tokens)
	util.CheckEqualStringer(t, "3.5", <-tokens)
	util.CheckEqualStringer(t, ")", <-tokens)
	util.CheckEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester3", "\\a")
	util.CheckEqualStringer(t, "\\a", <-tokens)

	fmt.Printf("START\n")
	_, tokens = Scan("tester4", "(list \\a)")
	util.CheckEqualStringer(t, "(", <-tokens)
	util.CheckEqualStringer(t, "list", <-tokens)
	util.CheckEqualStringer(t, "\\a", <-tokens)
	util.CheckEqualStringer(t, ")", <-tokens)
	fmt.Printf("END\n")
}
*/

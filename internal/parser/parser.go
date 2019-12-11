package parser

import (
	"strconv"
	"strings"

	"github.com/onlyafly/oakblue/internal/cst"
	"github.com/onlyafly/oakblue/internal/syntax"
	"github.com/onlyafly/oakblue/internal/util"
)

// Parse accepts a string and the name of the source of the code, and returns
// the Oakblue nodes therein, along with a list of any errors found.
func Parse(input string, sourceName string, errorList *syntax.ErrorList) (cst.Listing, error) {
	s, _ := Scan(sourceName, input)
	s.errorHandler = func(t Token, message string) {
		errorList.Add(t.Loc, message)
	}

	p := &parser{s: s}
	lines := parseLines(p, errorList)

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return cst.Listing(lines), nil
}

////////// Parser

type parser struct {
	s              *Scanner
	lookahead      [2]Token // two-token lookahead
	lookaheadCount int
}

func (p *parser) next() Token {
	if p.lookaheadCount > 0 {
		p.lookaheadCount--
	} else {
		p.lookahead[0] = <-p.s.Tokens
	}
	return p.lookahead[p.lookaheadCount]
}

/* TODO still needed?
func (p *parser) backup() {
	p.lookaheadCount++
}
*/

func (p *parser) peek() Token {
	if p.lookaheadCount > 0 {
		return p.lookahead[p.lookaheadCount-1]
	}

	p.lookaheadCount = 1
	p.lookahead[0] = <-p.s.Tokens
	return p.lookahead[0]
}

func (p *parser) inputEmpty() bool {
	c := p.peek().Code
	if c == TcEOF || c == TcError {
		return true
	}

	return false
}

func (p *parser) skipEmptyLines() {
	for p.peek().Code == TcNewline {
		p.next()
	}
}

////////// Parsing

func parseLines(p *parser, errors *syntax.ErrorList) []*cst.Line {
	var lines []*cst.Line

	p.skipEmptyLines()

	for !p.inputEmpty() {
		lines = append(lines, parseLine(p, errors))
	}
	return lines
}

func parseLine(p *parser, errors *syntax.ErrorList) *cst.Line {
	var nodes []cst.Node

	for !p.inputEmpty() {
		if p.peek().Code == TcNewline {
			p.next()
			break
		}
		nodes = append(nodes, parseNode(p, errors))
	}
	return cst.NewLine(nodes)
}

func parseNode(p *parser, errors *syntax.ErrorList) cst.Node {
	token := p.next()

	switch token.Code {
	case TcError:
		errors.Add(token.Loc, "Error token: "+token.String())
	case TcLeftParen:
		/* TODO delete
		var list []cst.Node
		for p.peek().Code != TcRightParen {
			if p.peek().Code == TcEOF || p.peek().Code == TcError {
				errors.Add(token.Loc, "Unbalanced parentheses")
				p.next()
				return &cst.Invalid{Location: token.Loc}
			}
			list = append(list, parseNode(p, errors))
		}
		p.next()
		return &cst.List{Nodes: list, Location: token.Loc}
		*/
	case TcRightParen:
		errors.Add(token.Loc, "Unbalanced parentheses")
	case TcNumber:
		return parseInteger(token, errors)
	case TcHex:
		return parseHex(token, errors)
	case TcSymbol:
		return parseSymbol(token, errors)
	case TcString:
		return parseString(token, errors)
	case TcChar:
		//TODO delete: return parseChar(token, errors)
	case TcSingleQuote:
		return parseQuote(p, errors)
	default:
		errors.Add(token.Loc, "Unrecognized token: "+token.String())
	}

	return &cst.Invalid{Location: token.Loc}
}

func parseQuote(p *parser, errors *syntax.ErrorList) cst.Node {
	/* TODO remove
	node := parseNode(p, errors)
	var list []cst.Node
	list = append(list, &cst.Symbol{Name: "quote"}, node)
	return &cst.List{Nodes: list}
	*/
	return &cst.Invalid{}
}

func parseInteger(t Token, errors *syntax.ErrorList) *cst.Integer {
	x, err := strconv.ParseInt(t.Value, 10, 32)

	if err != nil {
		errors.Add(t.Loc, "Invalid integer: "+t.Value)
		return &cst.Integer{Value: 0, Location: t.Loc}
	}

	return &cst.Integer{Value: int(x), Location: t.Loc}
}

func parseHex(t Token, errors *syntax.ErrorList) *cst.Hex {
	hexString := t.Value
	if strings.HasPrefix(hexString, "0x") {
		hexString = util.TrimLeftChars(hexString, 2)
	} else if strings.HasPrefix(hexString, "x") {
		hexString = util.TrimLeftChars(hexString, 1)
	}

	x, err := strconv.ParseUint(hexString, 16, 16)
	if err != nil {
		errors.Add(t.Loc, "Invalid hex: "+t.Value)
		return &cst.Hex{Value: uint16(x), Location: t.Loc}
	}

	return &cst.Hex{Value: uint16(x), Location: t.Loc}
}

func parseSymbol(t Token, errors *syntax.ErrorList) cst.Node {
	if t.Value == "nil" {
		return &cst.Invalid{Location: t.Loc}
	}
	return &cst.Symbol{Name: t.Value, Location: t.Loc}
}

func parseString(t Token, errors *syntax.ErrorList) *cst.Str {
	content := t.Value[1 : len(t.Value)-1]
	return &cst.Str{Value: content, Location: t.Loc}
}

////////// Helper Procedures

/* TODO still needed?
func ensureSymbol(n cst.Node) *cst.Symbol {
	if v, ok := n.(*cst.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}
*/

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
		errorList.Add(t, message)
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

/* TODO: still needed?
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
		errors.Add(token, "Error token: "+token.String())
	case TcLeftParen:
		/* TODO: delete
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
		errors.Add(token, "Unbalanced parentheses")
	case TcDecimalNumber:
		return parseDecimalNumber(token, errors)
	case TcHexNumber:
		return parseHexNumber(token, errors)
	case TcSymbol:
		if p.peek().Code == TcColon {
			p.next()
			return &cst.Label{Name: token.Value, Location: token.Location}
		}
		return parseSymbol(token, errors)
	case TcString:
		return parseString(token, errors)
	case TcChar:
		//TODO: delete this: return parseChar(token, errors)
		panic("char literal not implemented")
	case TcColon:
		errors.Add(token, "Colon expected only after symbol")
	case TcSingleQuote:
		return parseQuote(p, errors)
	default:
		errors.Add(token, "Unrecognized token: "+token.String())
	}

	return &cst.Invalid{Location: token.Location}
}

func parseQuote(p *parser, errors *syntax.ErrorList) cst.Node {
	/* TODO: remove
	node := parseNode(p, errors)
	var list []cst.Node
	list = append(list, &cst.Symbol{Name: "quote"}, node)
	return &cst.List{Nodes: list}
	*/
	return &cst.Invalid{}
}

func parseDecimalNumber(t Token, errors *syntax.ErrorList) *cst.DecimalNumber {
	intString := t.Value
	if strings.HasPrefix(t.Value, "#") {
		intString = util.TrimLeftChars(t.Value, 1)
	}

	x, err := strconv.ParseInt(intString, 10, 32)

	if err != nil {
		errors.Add(t, "Invalid integer: "+t.Value)
		return &cst.DecimalNumber{Value: 0, Location: t.Location}
	}

	return &cst.DecimalNumber{Value: int(x), Location: t.Location}
}

func parseHexNumber(t Token, errors *syntax.ErrorList) *cst.HexNumber {
	hexString := t.Value
	if strings.HasPrefix(hexString, "0x") {
		hexString = util.TrimLeftChars(hexString, 2)
	} else if strings.HasPrefix(hexString, "x") {
		hexString = util.TrimLeftChars(hexString, 1)
	}

	x, err := strconv.ParseUint(hexString, 16, 16)
	if err != nil {
		errors.Add(t, "Invalid hex: "+t.Value)
		return &cst.HexNumber{Value: uint16(x), Location: t.Location}
	}

	return &cst.HexNumber{Value: uint16(x), Location: t.Location}
}

func parseSymbol(t Token, errors *syntax.ErrorList) cst.Node {
	return &cst.Symbol{Name: t.Value, Location: t.Location}
}

func parseString(t Token, errors *syntax.ErrorList) *cst.Str {
	content := t.Value[1 : len(t.Value)-1]
	return &cst.Str{Value: content, Location: t.Location}
}

////////// Helper Procedures

/* TODO: still needed?
func ensureSymbol(n cst.Node) *cst.Symbol {
	if v, ok := n.(*cst.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}
*/

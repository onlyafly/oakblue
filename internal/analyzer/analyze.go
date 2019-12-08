package analyzer

/*
import (
	"strconv"

	"github.com/onlyafly/oakblue/internal/cst"
	"github.com/onlyafly/oakblue/internal/syntax"
)
*/

/*
func Analyze(input *cst.Listing) (ast.Listing, error) {
	s, _ := Scan(sourceName, input)
	errorList := Newsyntax.ErrorList()
	s.errorHandler = func(t Token, message string) {
		errorList.Add(t.Loc, message)
	}

	p := &parser{s: s}
	statements := parseStatements(p, errorList)

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return cst.Listing(statements), nil
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

func parseStatements(p *parser, errors *syntax.ErrorList) []*cst.Statement {
	var statements []*cst.Statement

	p.skipEmptyLines()

	for !p.inputEmpty() {
		statements = append(statements, parseStatement(p, errors))
	}
	return statements
}

func parseStatement(p *parser, errors *syntax.ErrorList) *cst.Statement {
	var nodes []cst.Node

	for !p.inputEmpty() {
		if p.peek().Code == TcNewline {
			p.next()
			break
		}
		nodes = append(nodes, parseNode(p, errors))
	}
	return cst.NewStatement(nodes)
}

func parseNode(p *parser, errors *syntax.ErrorList) cst.Node {
	token := p.next()

	switch token.Code {
	case TcError:
		errors.Add(token.Loc, "Error token: "+token.String())
	case TcRightParen:
		errors.Add(token.Loc, "Unbalanced parentheses")
	case TcNumber:
		return parseInteger(token, errors)
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

func ensureSymbol(n cst.Node) *cst.Symbol {
	if v, ok := n.(*cst.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}
*/

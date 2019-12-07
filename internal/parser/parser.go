package parser

import (
	"strconv"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/spec"
)

// Parse accepts a string and the name of the source of the code, and returns
// the Oakblue nodes therein, along with a list of any errors found.
func Parse(input string, sourceName string) (ast.Program, ParserErrorList) {
	/* TODO readd this
	s, _ := Scan(sourceName, input)
	errorList := NewParserErrorList()
	s.errorHandler = func(t Token, message string) {
		errorList.Add(t.Loc, message)
	}

	p := &parser{s: s}
	nodes := parseNodes(p, &errorList)

	if errorList.Len() > 0 {
		return nil, errorList
	}
	*/

	// TODO remove this
	var program ast.Program
	program = append(program, ast.NewStatement([]ast.Node{ast.NewOp(spec.OP_ADD)}))

	return program, nil
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

////////// Parsing

func parseNodes(p *parser, errors *ParserErrorList) []ast.Node {
	var nodes []ast.Node
	for !p.inputEmpty() {
		nodes = append(nodes, parseNode(p, errors))
	}
	return nodes
}

func parseNode(p *parser, errors *ParserErrorList) ast.Node {
	token := p.next()

	switch token.Code {
	case TcError:
		errors.Add(token.Loc, "Error token: "+token.String())
	case TcLeftParen:
		/* TODO delete
		var list []ast.Node
		for p.peek().Code != TcRightParen {
			if p.peek().Code == TcEOF || p.peek().Code == TcError {
				errors.Add(token.Loc, "Unbalanced parentheses")
				p.next()
				return &ast.Invalid{Location: token.Loc}
			}
			list = append(list, parseNode(p, errors))
		}
		p.next()
		return &ast.List{Nodes: list, Location: token.Loc}
		*/
	case TcRightParen:
		errors.Add(token.Loc, "Unbalanced parentheses")
	case TcNumber:
		return parseNumber(token, errors)
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

	return &ast.Invalid{Location: token.Loc}
}

func parseQuote(p *parser, errors *ParserErrorList) ast.Node {
	/* TODO remove
	node := parseNode(p, errors)
	var list []ast.Node
	list = append(list, &ast.Symbol{Name: "quote"}, node)
	return &ast.List{Nodes: list}
	*/
	return &ast.Invalid{}
}

func parseNumber(t Token, errors *ParserErrorList) *ast.Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)

	if ferr != nil {
		errors.Add(t.Loc, "Invalid number: "+t.Value)
		return &ast.Number{Value: 0.0, Location: t.Loc}
	}

	return &ast.Number{Value: f, Location: t.Loc}
}

func parseSymbol(t Token, errors *ParserErrorList) ast.Node {
	if t.Value == "nil" {
		return &ast.Invalid{Location: t.Loc}
	}
	return &ast.Symbol{Name: t.Value, Location: t.Loc}
}

func parseString(t Token, errors *ParserErrorList) *ast.Str {
	content := t.Value[1 : len(t.Value)-1]
	return &ast.Str{Value: content, Location: t.Loc}
}

////////// Helper Procedures

/* TODO still needed?
func ensureSymbol(n ast.Node) *ast.Symbol {
	if v, ok := n.(*ast.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}
*/

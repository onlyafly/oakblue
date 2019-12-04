package parser

import "github.com/onlyafly/oakblue/internal/ast"

// Parse a string into nodes
func Parse(input string, sourceName string) (ast.Nodes, SyntaxErrorList) {
	var nodes []ast.Node
	nodes = append(nodes, ast.NewStr("asdf"))
	return nodes, nil
}

package interpreter

import (
	"github.com/onlyafly/oakblue/internal/ast"
	"io"
)

//TODO var writer io.Writer
//TODO var readLine func() string

// Eval a node
func Eval(e Env, n ast.Node, w io.Writer, rl func() string) (result ast.Node, err error) {
	return ast.NewStr("ert"), nil
}

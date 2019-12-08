package interpreter

import (
	"github.com/onlyafly/oakblue/internal/cst"
	"io"
)

//TODO var writer io.Writer
//TODO var readLine func() string

// Eval a node
func Eval(e Env, program cst.Listing, w io.Writer, rl func() string) (result cst.Node, err error) {
	return cst.NewStr("ert"), nil
}

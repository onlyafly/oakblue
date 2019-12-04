package interpreter

import (
	"fmt"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/token"
)

// EvalError represents an error that occurs during evaluation.
type EvalError struct {
	SuperMessage string
	Message      string
	location     *token.Location
}

// NewEvalError returns a new EvalError
func NewEvalError(superMessage, message string, location *token.Location) *EvalError {
	return &EvalError{superMessage, message, location}
}

// Implements the error interface
func (e *EvalError) Error() string {
	if e.location != nil {
		return fmt.Sprintf("%v (%v: %v): %v", e.SuperMessage, e.location.Filename, e.location.Line, e.Message)
	}

	return fmt.Sprintf("%v: %v", e.SuperMessage, e.Message)
}

func panicEvalError(n ast.Node, s string) {
	var loc *token.Location
	if n != nil {
		loc = n.Loc()
	}
	panic(NewEvalError("Evaluation error", s, loc))
}

/* TODO
func panicApplicationError(n ast.Node, s string) {
	var loc *token.Location
	if n != nil {
		loc = n.Loc()
	}
	panic(NewEvalError("Application panic", s, loc))
}
*/

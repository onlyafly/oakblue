package parser

import (
	"fmt"
	"github.com/onlyafly/oakblue/internal/token"
)

////////// ParserError

// SyntaxError describes a parsing error
type SyntaxError struct {
	Loc     *token.Location
	Message string
}

// Implements the error interface
func (pe *SyntaxError) Error() string {
	if pe.Loc != nil {
		return fmt.Sprintf("Parsing error (%v: %v): %v", pe.Loc.Filename, pe.Loc.Line, pe.Message)
	}

	return fmt.Sprintf("Parsing error: %v", pe.Message)
}

////////// ParserErrorList

// SyntaxErrorList is a list of ParserError pointers; Implements the error interface.
type SyntaxErrorList []*SyntaxError

// NewSyntaxErrorList creates a new syntax error list
func NewSyntaxErrorList() SyntaxErrorList {
	return make(SyntaxErrorList, 0)
}

// Add an error
func (p *SyntaxErrorList) Add(loc *token.Location, msg string) {
	*p = append(*p, &SyntaxError{loc, msg})
}

func (p SyntaxErrorList) Error() string {
	return p.String()
}

// Len of the list
func (p SyntaxErrorList) Len() int {
	return len(p)
}

func (p SyntaxErrorList) String() (s string) {
	for i, e := range p {
		s += e.Error()

		if i != len(p)-1 {
			s += "\n"
		}
	}

	return s
}

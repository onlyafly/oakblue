package parser

import (
	"fmt"

	"github.com/onlyafly/oakblue/internal/token"
)

////////// ParserError

type ParserError struct {
	Loc     *token.Location
	Message string
}

// Implements the error interface
func (pe *ParserError) Error() string {
	if pe.Loc != nil {
		return fmt.Sprintf("Parsing error (%v: %v): %v", pe.Loc.Filename, pe.Loc.Line, pe.Message)
	}

	return fmt.Sprintf("Parsing error: %v", pe.Message)
}

////////// ParserErrorList

// ParserErrorList is a list of ParserError pointers.
// Implements the error interface.
type ParserErrorList struct {
	Errors []*ParserError
}

func NewParserErrorList() *ParserErrorList {
	return &ParserErrorList{Errors: make([]*ParserError, 0)}
}

func (p *ParserErrorList) Add(loc *token.Location, msg string) {
	p.Errors = append(p.Errors, &ParserError{loc, msg})
}

func (p ParserErrorList) Error() string {
	return p.String()
}

func (p ParserErrorList) Len() int {
	return len(p.Errors)
}

func (p ParserErrorList) String() (s string) {
	for i, e := range p.Errors {
		s += e.Error()

		if i != len(p.Errors)-1 {
			s += "\n"
		}
	}

	return s
}

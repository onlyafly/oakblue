package syntax

import (
	"fmt"
)

type Error struct {
	Loc     *Location
	Message string
}

// Implements the error interface
func (pe *Error) Error() string {
	if pe.Loc != nil {
		return fmt.Sprintf("Parsing error (%v: %v): %v", pe.Loc.Filename, pe.Loc.Line, pe.Message)
	}

	return fmt.Sprintf("Parsing error: %v", pe.Message)
}

// ErrorList is a list of Error pointers.
// Implements the error interface.
type ErrorList struct {
	Errors []*Error
}

func NewErrorList() *ErrorList {
	return &ErrorList{Errors: make([]*Error, 0)}
}

func (p *ErrorList) Add(loc *Location, msg string) {
	p.Errors = append(p.Errors, &Error{loc, msg})
}

func (p ErrorList) Error() string {
	return p.String()
}

func (p ErrorList) Len() int {
	return len(p.Errors)
}

func (p ErrorList) String() (s string) {
	for i, e := range p.Errors {
		s += e.Error()

		if i != len(p.Errors)-1 {
			s += "\n"
		}
	}

	return s
}

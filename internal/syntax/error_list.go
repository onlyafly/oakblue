package syntax

import (
	"fmt"
)

type Error struct {
	Loc     *Location
	Message string
	Kind    string
}

// Implements the error interface
func (e *Error) Error() string {
	if e.Loc != nil {
		return fmt.Sprintf("%s error (%v: %v): %v", e.Kind, e.Loc.Filename, e.Loc.Line, e.Message)
	}

	return fmt.Sprintf("%s error: %v", e.Kind, e.Message)
}

// ErrorList is a list of Error pointers.
// Implements the error interface.
type ErrorList struct {
	Errors []*Error
	Kind   string
}

func NewErrorList(kind string) *ErrorList {
	return &ErrorList{Errors: make([]*Error, 0), Kind: kind}
}

func (el *ErrorList) Add(l HasLocation, msg string) {
	el.Errors = append(el.Errors, &Error{l.Loc(), msg, el.Kind})
}

func (el ErrorList) Error() string {
	return el.String()
}

func (el ErrorList) Len() int {
	return len(el.Errors)
}

func (el ErrorList) String() (s string) {
	for i, e := range el.Errors {
		s += e.Error()

		if i != len(el.Errors)-1 {
			s += "\n"
		}
	}

	return s
}

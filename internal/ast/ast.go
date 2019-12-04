package ast

import (
	"fmt"

	"github.com/onlyafly/oakblue/internal/token"
)

// Nodes type represents an array of Nodes.
type Nodes []Node

// Node represents a parsed node.
type Node interface {
	fmt.Stringer
	FriendlyStringer

	Equals(Node) bool
	TypeName() string
	Loc() *token.Location
}

// Str is a node
type Str struct {
	Value      string
	annotation Node
	Location   *token.Location
}

func NewStr(value string) *Str { return &Str{Value: value} }

func (s *Str) String() string         { return "\"" + s.Value + "\"" }
func (s *Str) FriendlyString() string { return s.Value }

//TODO func (s *Str) isExpr() bool           { return true }

func (s *Str) Annotation() Node     { return s.annotation }
func (s *Str) SetAnnotation(n Node) { s.annotation = n }
func (s *Str) Equals(n Node) bool   { return s.Value == asStr(n).Value }
func (s *Str) TypeName() string     { return "string" }
func (s *Str) Loc() *token.Location { return s.Location }

func asStr(n Node) *Str {
	if result, ok := n.(*Str); ok {
		return result
	}
	return &Str{}
}

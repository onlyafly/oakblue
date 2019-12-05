package ast

import (
	"fmt"

	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/token"
)

type Program []*Statement

type Statement struct {
	Nodes []Node
}

func NewStatement(nodes []Node) *Statement { return &Statement{Nodes: nodes} }

// Node represents a parsed node.
type Node interface {
	fmt.Stringer
	Loc() *token.Location
}

type Op struct {
	Opcode   int
	Location *token.Location
}

func NewOp(opcode int) *Op         { return &Op{Opcode: opcode} }
func (o *Op) String() string       { return spec.OpcodeNames[o.Opcode] }
func (o *Op) Loc() *token.Location { return o.Location }

// Str is a node
type Str struct {
	Value    string
	Location *token.Location
}

func NewStr(value string) *Str      { return &Str{Value: value} }
func (s *Str) String() string       { return "\"" + s.Value + "\"" }
func (s *Str) Loc() *token.Location { return s.Location }

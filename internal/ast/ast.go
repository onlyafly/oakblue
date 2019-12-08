package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/token"
)

type Program []*Statement

func (p *Program) String() string {
	return strings.Join(statementsToStrings(*p), "\n")
}

type Statement struct {
	Nodes []Node
}

func statementsToStrings(statements []*Statement) []string {
	return statementsToStringsWithFunc(statements, func(x *Statement) string { return x.String() })
}
func statementsToStringsWithFunc(statements []*Statement, convert func(x *Statement) string) []string {
	strings := make([]string, len(statements))
	for i, x := range statements {
		strings[i] = convert(x)
	}
	return strings
}

func NewStatement(nodes []Node) *Statement { return &Statement{Nodes: nodes} }
func (x *Statement) String() string {
	return strings.Join(nodesToStrings(x.Nodes), " ")
}

// Node represents a parsed node.
type Node interface {
	fmt.Stringer
	Loc() *token.Location
}

func nodesToStrings(nodes []Node) []string {
	return nodesToStringsWithFunc(nodes, func(n Node) string { return n.String() })
}
func nodesToStringsWithFunc(nodes []Node, convert func(n Node) string) []string {
	strings := make([]string, len(nodes))
	for i, node := range nodes {
		strings[i] = convert(node)
	}
	return strings
}

type Op struct {
	Opcode   int
	Location *token.Location
}

func NewOp(opcode int) *Op         { return &Op{Opcode: opcode} }
func (o *Op) String() string       { return spec.OpcodeNames[o.Opcode] }
func (o *Op) Loc() *token.Location { return o.Location }

// Symbol is a node
type Symbol struct {
	Name     string
	Location *token.Location
}

func NewSymbol(name string) *Symbol    { return &Symbol{Name: name} }
func (s *Symbol) String() string       { return s.Name }
func (s *Symbol) Loc() *token.Location { return s.Location }

// Str is a node
type Str struct {
	Value    string
	Location *token.Location
}

func NewStr(value string) *Str      { return &Str{Value: value} }
func (s *Str) String() string       { return "\"" + s.Value + "\"" }
func (s *Str) Loc() *token.Location { return s.Location }

type Integer struct {
	Value    int
	Location *token.Location
}

func NewInteger(value int) *Integer { return &Integer{Value: value} }
func (x *Integer) String() string {
	rep := strconv.FormatInt(int64(x.Value), 10)
	return rep
}
func (x *Integer) Loc() *token.Location { return x.Location }

type Invalid struct {
	Value    string
	Location *token.Location
}

func (x *Invalid) String() string       { return "INVALID<" + x.Value + ">" }
func (x *Invalid) Loc() *token.Location { return x.Location }

package cst

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onlyafly/oakblue/internal/syntax"
)

type Listing []*Statement

func (l *Listing) String() string {
	return strings.Join(statementsToStrings(*l), "\n")
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
	Loc() *syntax.Location
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

// Symbol is a node
type Symbol struct {
	Name     string
	Location *syntax.Location
}

func NewSymbol(name string) *Symbol     { return &Symbol{Name: name} }
func (s *Symbol) String() string        { return s.Name }
func (s *Symbol) Loc() *syntax.Location { return s.Location }

// Str is a node
type Str struct {
	Value    string
	Location *syntax.Location
}

func NewStr(value string) *Str       { return &Str{Value: value} }
func (s *Str) String() string        { return "\"" + s.Value + "\"" }
func (s *Str) Loc() *syntax.Location { return s.Location }

type Integer struct {
	Value    int
	Location *syntax.Location
}

func NewInteger(value int) *Integer { return &Integer{Value: value} }
func (x *Integer) String() string {
	rep := strconv.FormatInt(int64(x.Value), 10)
	return rep
}
func (x *Integer) Loc() *syntax.Location { return x.Location }

type Invalid struct {
	Value    string
	Location *syntax.Location
}

func (x *Invalid) String() string        { return "INVALID<" + x.Value + ">" }
func (x *Invalid) Loc() *syntax.Location { return x.Location }

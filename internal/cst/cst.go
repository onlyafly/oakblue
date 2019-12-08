package cst

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/onlyafly/oakblue/internal/syntax"
)

type Listing []*Line

func (l *Listing) String() string {
	return strings.Join(linesToStrings(*l), "\n")
}

type Line struct {
	Nodes []Node
}

func linesToStrings(lines []*Line) []string {
	return linesToStringsWithFunc(lines, func(x *Line) string { return x.String() })
}
func linesToStringsWithFunc(lines []*Line, convert func(x *Line) string) []string {
	strings := make([]string, len(lines))
	for i, x := range lines {
		strings[i] = convert(x)
	}
	return strings
}

func NewLine(nodes []Node) *Line { return &Line{Nodes: nodes} }
func (x *Line) String() string {
	return strings.Join(nodesToStrings(x.Nodes), " ")
}
func (x *Line) Loc() *syntax.Location {
	if len(x.Nodes) > 0 {
		return x.Nodes[0].Loc()
	}
	return &syntax.Location{}
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

package ast

import (
	"fmt"
	"strings"

	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
)

type Program struct {
	Statements []Statement
	Symtab     *SymbolTable
	Origin     uint16
}

func NewProgram(xs []Statement, symtab *SymbolTable, origin uint16) *Program {
	return &Program{Statements: xs, Symtab: symtab, Origin: origin}
}

func (p *Program) String() string {
	return strings.Join(statementsToStrings(p.Statements), "\n")
}

type Statement interface {
	fmt.Stringer
	Loc() *syntax.Location
}

func statementsToStrings(statements []Statement) []string {
	return statementsToStringsWithFunc(statements, func(x Statement) string { return x.String() })
}
func statementsToStringsWithFunc(statements []Statement, convert func(x Statement) string) []string {
	strings := make([]string, len(statements))
	for i, x := range statements {
		strings[i] = convert(x)
	}
	return strings
}

type Instruction struct {
	Opcode      int
	Dr          int
	Sr1         int
	Sr2         int
	Mode        int
	Imm5        int
	Trapvect8   uint8
	PCOffset9   int
	Label       string
	BranchFlags *BranchFlags
	Location    *syntax.Location
}

func (x *Instruction) String() string {
	switch x.Opcode {
	case spec.OP_ADD:
		switch x.Mode {
		case 0:
			return fmt.Sprintf("ADD %s %s %s", spec.RegisterNames[x.Dr], spec.RegisterNames[x.Sr1], spec.RegisterNames[x.Sr2])
		case 1:
			return fmt.Sprintf("ADD %s %s %v", spec.RegisterNames[x.Dr], spec.RegisterNames[x.Sr1], x.Imm5)
		}
	default:
		return fmt.Sprintf("<UNRECOGNIZED OPCODE=%s>", spec.OpcodeNames[x.Opcode])
	}

	return fmt.Sprintf("<MALFORMED INSTRUCTION %#v>", x)
}

func (x *Instruction) Loc() *syntax.Location { return x.Location }

type InvalidStatement struct {
	MoreInformation string
	Location        *syntax.Location
}

func (x *InvalidStatement) String() string        { return "<INVALID STATEMENT: " + x.MoreInformation + ">" }
func (x *InvalidStatement) Loc() *syntax.Location { return x.Location }

type FillDirective struct {
	Value    uint16
	Location *syntax.Location
}

func (x *FillDirective) String() string        { return fmt.Sprintf(".FILL %d", x.Value) }
func (x *FillDirective) Loc() *syntax.Location { return x.Location }

type BranchFlags struct {
	N int
	Z int
	P int
}

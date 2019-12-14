package analyzer

import (
	"fmt"
	"strings"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/cst"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
)

func Analyze(input cst.Listing, errorList *syntax.ErrorList) (*ast.Program, error) {
	a := &analyzer{errors: errorList}
	statements := a.analyzeStatements(input)

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return ast.NewProgram(statements), nil
}

type analyzer struct {
	errors *syntax.ErrorList
}

func (a *analyzer) analyzeStatements(l cst.Listing) []ast.Statement {
	var statements []ast.Statement

	for _, line := range l {
		statements = append(statements, a.analyzeStatement(line))
	}

	return statements
}

func (a *analyzer) analyzeStatement(l *cst.Line) ast.Statement {
	firstNode := l.Nodes[0]

	switch v := firstNode.(type) {
	case *cst.Symbol:
		switch strings.ToUpper(v.Name) {
		case "ADD":
			return a.analyzeAddInstruction(l)
		case "AND":
			return a.analyzeAndInstruction(l)
		case "NOT":
			return a.analyzeNotInstruction(l)
		case "TRAP":
			return a.analyzeTrapInstruction(l)
		case ".FILL":
			return a.analyzeFillDirective(l)
		default:
			a.errors.Add(v.Loc(), "unrecognized operation name: "+v.Name)
		}
	default:
		a.errors.Add(v.Loc(), "unrecognized statement syntax")
	}

	return &ast.InvalidStatement{Location: firstNode.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeAddInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 3) {
		return &ast.InvalidStatement{}
	}

	dr := a.analyzeRegister(l.Nodes[1])
	sr1 := a.analyzeRegister(l.Nodes[2])

	switch arg3 := l.Nodes[3].(type) {
	case *cst.Symbol:
		sr2 := a.analyzeRegister(arg3)
		return &ast.Instruction{
			Opcode:   spec.OP_ADD,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     0,
			Sr2:      sr2,
			Location: l.Loc(),
		}
	case *cst.Integer:
		return &ast.Instruction{
			Opcode:   spec.OP_ADD,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     1,
			Imm5:     arg3.Value,
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg3.Loc(), "expected register or integer, got: "+arg3.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeAndInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 3) {
		return &ast.InvalidStatement{}
	}

	dr := a.analyzeRegister(l.Nodes[1])
	sr1 := a.analyzeRegister(l.Nodes[2])

	switch arg3 := l.Nodes[3].(type) {
	case *cst.Symbol:
		sr2 := a.analyzeRegister(arg3)
		return &ast.Instruction{
			Opcode:   spec.OP_AND,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     0,
			Sr2:      sr2,
			Location: l.Loc(),
		}
	case *cst.Integer:
		return &ast.Instruction{
			Opcode:   spec.OP_AND,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     1,
			Imm5:     arg3.Value,
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg3.Loc(), "expected register or integer, got: "+arg3.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeNotInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 2) {
		return &ast.InvalidStatement{}
	}

	dr := a.analyzeRegister(l.Nodes[1])
	sr := a.analyzeRegister(l.Nodes[2])

	return &ast.Instruction{
		Opcode:   spec.OP_NOT,
		Dr:       dr,
		Sr1:      sr,
		Location: l.Loc(),
	}
}

func (a *analyzer) analyzeTrapInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 1) {
		return &ast.InvalidStatement{}
	}

	switch arg := l.Nodes[1].(type) {
	case *cst.Hex:
		return &ast.Instruction{
			Opcode:    spec.OP_TRAP,
			Trapvect8: uint8(arg.Value),
			Location:  l.Loc(),
		}
	default:
		a.errors.Add(arg.Loc(), "expected hex, got: "+arg.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeFillDirective(l *cst.Line) ast.Statement {
	switch arg := l.Nodes[1].(type) {
	case *cst.Integer:
		return &ast.FillDirective{
			Value:    uint16(arg.Value),
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg.Loc(), "expected integer, got: "+arg.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeRegister(n cst.Node) int {
	switch v := n.(type) {
	case *cst.Symbol:
		switch v.Name {
		case "R0":
			return spec.R_R0
		case "R1":
			return spec.R_R1
		case "R2":
			return spec.R_R2
		case "R3":
			return spec.R_R3
		case "R4":
			return spec.R_R4
		case "R5":
			return spec.R_R5
		case "R6":
			return spec.R_R6
		case "R7":
			return spec.R_R7
		default:
			a.errors.Add(v.Loc(), "expected register, got: "+v.Name)
		}
	default:
		a.errors.Add(v.Loc(), "expected symbol, got: "+v.String())
	}

	return 0
}

func (a *analyzer) ensureLineArgs(l *cst.Line, argCount int) bool {
	if len(l.Nodes) != argCount+1 {
		a.errors.Add(l.Loc(), fmt.Sprintf("expected %d arguments, got: %d", argCount, len(l.Nodes)-1))
		return false
	}
	return true
}

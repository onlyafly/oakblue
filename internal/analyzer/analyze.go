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
	symtab := ast.NewSymbolTable()
	a := &analyzer{
		errors: errorList,
		symtab: symtab,
	}
	statements := a.analyzeStatements(input)

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return ast.NewProgram(statements, symtab), nil
}

type analyzer struct {
	errors *syntax.ErrorList
	symtab *ast.SymbolTable
}

func (a *analyzer) analyzeStatements(l cst.Listing) []ast.Statement {
	var statements []ast.Statement

	for i, line := range l {
		statements = append(statements, a.analyzeStatement(i, line))
	}

	return statements
}

func (a *analyzer) analyzeStatement(lineIndex int, l *cst.Line) ast.Statement {
	firstNode := l.Nodes[0]

	// Analyze the optional label
	switch v := firstNode.(type) {
	case *cst.Label:
		err := a.symtab.Insert(v.Name, uint16(lineIndex))
		if err != nil {
			a.errors.Add(v, "label redefined: "+v.String())
		}

		l = cst.NewLine(l.Nodes[1:])
		firstNode = l.Nodes[0]
	default:
		// Do nothing
	}

	switch v := firstNode.(type) {
	case *cst.Symbol:
		switch strings.ToUpper(v.Name) {
		case "ADD":
			return a.analyzeAddInstruction(l)
		case "AND":
			return a.analyzeAndInstruction(l)
		case "LD":
			return a.analyzeLdInstruction(l)
		case "NOT":
			return a.analyzeNotInstruction(l)
		case "TRAP":
			return a.analyzeTrapInstruction(l)
		case "HALT":
			return a.analyzeHaltInstruction(l)
		case ".FILL":
			return a.analyzeFillDirective(l)
		default:
			a.errors.Add(v, "unrecognized operation name: "+v.Name)
		}
	default:
		a.errors.Add(v, fmt.Sprintf("unrecognized statement syntax: %v", l))
	}

	return &ast.InvalidStatement{Location: firstNode.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeAddInstruction(l *cst.Line) ast.Statement {
	// TODO refactor the line args check out
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
		a.errors.Add(arg3, "expected register or integer, got: "+arg3.String())
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
		a.errors.Add(arg3, "expected register or integer, got: "+arg3.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeLdInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 2) {
		return &ast.InvalidStatement{}
	}

	dr := a.analyzeRegister(l.Nodes[1])

	switch arg2 := l.Nodes[2].(type) {
	case *cst.Symbol:
		sym := a.analyzeSymbol(arg2)
		return &ast.Instruction{
			Opcode:   spec.OP_LD,
			Dr:       dr,
			Label:    sym,
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg2, "expected symbol, got: "+arg2.String())
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
		a.errors.Add(arg, "expected hex, got: "+arg.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeHaltInstruction(l *cst.Line) ast.Statement {
	if !a.ensureLineArgs(l, 0) {
		return &ast.InvalidStatement{}
	}

	return &ast.Instruction{
		Opcode:    spec.OP_TRAP,
		Trapvect8: uint8(spec.TRAPVECT_HALT),
		Location:  l.Loc(),
	}
}

func (a *analyzer) analyzeFillDirective(l *cst.Line) ast.Statement {
	switch arg := l.Nodes[1].(type) {
	case *cst.Integer:
		return &ast.FillDirective{
			Value:    uint16(arg.Value),
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg, "expected integer, got: "+arg.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeRegister(n cst.Node) int {
	switch v := n.(type) {
	case *cst.Symbol:
		switch strings.ToUpper(v.Name) {
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
			a.errors.Add(v, "expected register, got: "+v.Name)
		}
	default:
		a.errors.Add(v, "expected symbol, got: "+v.String())
	}

	return 0
}

func (a *analyzer) analyzeSymbol(sym *cst.Symbol) string {
	return sym.Name
}

func (a *analyzer) ensureLineArgs(l *cst.Line, argCount int) bool {
	if len(l.Nodes) != argCount+1 {
		a.errors.Add(l, fmt.Sprintf("expected %d arguments, got: %d", argCount, len(l.Nodes)-1))
		return false
	}
	return true
}

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

	return ast.NewProgram(statements, symtab, a.customOrigin), nil
}

type analyzer struct {
	errors       *syntax.ErrorList
	symtab       *ast.SymbolTable
	customOrigin uint16
}

func (a *analyzer) analyzeStatements(l cst.Listing) []ast.Statement {
	var statements []ast.Statement

	var lineIndex uint16
	for _, line := range l {
		statement, statementSize := a.analyzeStatement(lineIndex, line)
		if statement != nil {
			statements = append(statements, statement)
		}
		lineIndex += statementSize
	}

	return statements
}

func (a *analyzer) analyzeStatement(lineIndex uint16, l *cst.Line) (ast.Statement, uint16) {
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
			return a.analyzeAddInstruction(l), 1
		case "AND":
			return a.analyzeAndInstruction(l), 1
		case "LD":
			return a.analyzeLdInstruction(l), 1
		case "NOT":
			return a.analyzeNotInstruction(l), 1
		case "TRAP":
			return a.analyzeTrapInstruction(l), 1
		case "HALT":
			return a.analyzeHaltInstruction(l), 1
		case ".FILL":
			return a.analyzeFillDirective(l), 1
		case ".ORIG":
			return a.analyzeOrigDirective(l, lineIndex), 0 // .ORIG directive has zero size
		default:
			a.errors.Add(v, "unrecognized operation name: "+v.Name)
		}
	default:
		a.errors.Add(v, fmt.Sprintf("unrecognized statement syntax: %v", l))
	}

	return &ast.InvalidStatement{Location: firstNode.Loc(), MoreInformation: l.String()}, 0
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
	default:
		imm5 := a.analyzeNumber(l.Nodes[3], "ADD")

		if imm5 >= 16 {
			a.errors.Add(l.Nodes[3], fmt.Sprintf("immediate argument to ADD must be less than 16: %d", imm5))
		}

		return &ast.Instruction{
			Opcode:   spec.OP_ADD,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     1,
			Imm5:     imm5,
			Location: l.Loc(),
		}
	}
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
	default:
		imm5 := a.analyzeNumber(l.Nodes[3], "AND")

		if imm5 >= 16 {
			a.errors.Add(l.Nodes[3], fmt.Sprintf("immediate argument to AND must be less than 16: %d", imm5))
		}

		return &ast.Instruction{
			Opcode:   spec.OP_AND,
			Dr:       dr,
			Sr1:      sr1,
			Mode:     1,
			Imm5:     imm5,
			Location: l.Loc(),
		}
	}
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

	trapvect8 := a.analyzeNumber(l.Nodes[1], "TRAP")

	return &ast.Instruction{
		Opcode:    spec.OP_TRAP,
		Trapvect8: uint8(trapvect8),
		Location:  l.Loc(),
	}
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
	case *cst.DecimalNumber:
		return &ast.FillDirective{
			Value:    uint16(arg.Value),
			Location: l.Loc(),
		}
	default:
		a.errors.Add(arg, "expected integer, got: "+arg.String())
	}

	return &ast.InvalidStatement{Location: l.Loc(), MoreInformation: l.String()}
}

func (a *analyzer) analyzeOrigDirective(l *cst.Line, lineIndex uint16) ast.Statement {
	if lineIndex != 0 {
		a.errors.Add(l, ".ORIG directive must appear before first instruction")
	}

	a.customOrigin = uint16(a.analyzeNumber(l.Nodes[1], ".ORIG"))

	return nil
}

func (a *analyzer) analyzeNumber(n cst.Node, instructionName string) int {
	switch x := n.(type) {
	case *cst.HexNumber:
		return int(x.Value)
	case *cst.DecimalNumber:
		return x.Value
	default:
		a.errors.Add(x, instructionName+" expected number, got: "+x.String())
		return 0
	}
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

package analyzer

import (
	"fmt"
	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/cst"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
)

func Analyze(input cst.Listing) (*ast.Program, error) {
	errorList := syntax.NewErrorList()

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
		switch v.Name {
		case "ADD":
			return a.analyzeAddInstruction(l)
		default:
			a.errors.Add(v.Loc(), "unrecognized operation name")
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
			Opcode: spec.OP_ADD,
			Dr:     dr,
			Sr1:    sr1,
			Mode:   0,
			Sr2:    sr2,
		}
	case *cst.Integer:
		return &ast.Instruction{
			Opcode: spec.OP_ADD,
			Dr:     dr,
			Sr1:    sr1,
			Mode:   1,
			Imm5:   arg3.Value,
		}
	default:
		a.errors.Add(arg3.Loc(), "expected register or integer, got: "+arg3.String())
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
			return spec.R_R0
		case "R2":
			return spec.R_R0
		case "R3":
			return spec.R_R0
		case "R4":
			return spec.R_R0
		case "R5":
			return spec.R_R0
		case "R6":
			return spec.R_R0
		case "R7":
			return spec.R_R0
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

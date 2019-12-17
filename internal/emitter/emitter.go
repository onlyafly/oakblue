package emitter

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/onlyafly/oakblue/internal/ast"
	"github.com/onlyafly/oakblue/internal/spec"
	"github.com/onlyafly/oakblue/internal/syntax"
)

// Emit emits an assembled binary image
func Emit(p *ast.Program, errorList *syntax.ErrorList) ([]byte, error) {
	var buf bytes.Buffer
	m := &emitter{errors: errorList, buf: &buf, tab: p.Symtab}

	if len(p.Statements) == 0 {
		return buf.Bytes(), nil
	}

	// Write header with origin
	if p.Origin == 0 {
		m.write(spec.DefaultOrigin, p.Statements[0])
	} else {
		m.write(p.Origin, p.Statements[0])
	}

	for pc, s := range p.Statements {
		switch v := s.(type) {
		case *ast.Instruction:
			m.emitInstruction(uint16(pc), v)
		case *ast.FillDirective:
			m.emitFillDirective(v)
		default:
			m.errors.Add(v, "unexpected statement type: "+v.String())
		}
	}

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return buf.Bytes(), nil
}

type emitter struct {
	errors *syntax.ErrorList
	buf    *bytes.Buffer
	tab    *ast.SymbolTable
}

func (m *emitter) emitInstruction(pc uint16, inst *ast.Instruction) {
	switch inst.Opcode {
	case spec.OP_ST, spec.OP_JSR,
		spec.OP_LDR, spec.OP_STR, spec.OP_RTI, spec.OP_LDI,
		spec.OP_STI, spec.OP_JMP, spec.OP_RES, spec.OP_LEA:
		m.errors.Add(inst, "emitter hasn't yet implemented this instruction: "+spec.OpcodeNames[inst.Opcode]) // TODO: implement these instructions
	case spec.OP_ADD:
		var x int
		x = spec.OP_ADD << 12
		x |= inst.Dr << 9
		x |= inst.Sr1 << 6

		switch inst.Mode {
		case 0:
			x |= 0 << 5
			x |= inst.Sr2
		case 1:
			x |= 1 << 5
			x |= inst.Imm5 & 0b11111
		default:
			m.errors.Add(inst, "unknown mode")
		}

		m.write(uint16(x), inst)
	case spec.OP_AND:
		var x int
		x = spec.OP_AND << 12
		x |= inst.Dr << 9
		x |= inst.Sr1 << 6

		switch inst.Mode {
		case 0:
			x |= 0 << 5
			x |= inst.Sr2
		case 1:
			x |= 1 << 5
			x |= inst.Imm5 & 0b11111
		default:
			m.errors.Add(inst, "unknown mode")
		}

		m.write(uint16(x), inst)
	case spec.OP_BR:

		var nzp int
		nzp |= (inst.BranchFlags.N & 0b1) << 2
		nzp |= (inst.BranchFlags.Z & 0b1) << 1
		nzp |= inst.BranchFlags.P & 0b1

		var x int
		x = spec.OP_BR << 12
		x |= (nzp & 0b111) << 9

		if len(inst.Label) != 0 {
			x |= m.labelToOffset(inst.Label, 0b111111111, pc, inst)
		} else {
			x |= inst.PCOffset9 & 0b111111111
		}

		m.write(uint16(x), inst)
	case spec.OP_LD:
		var x int
		x = spec.OP_LD << 12
		x |= inst.Dr << 9
		x |= m.labelToOffset(inst.Label, 0b111111111, pc, inst)

		m.write(uint16(x), inst)
	case spec.OP_NOT:
		var x int
		x = spec.OP_NOT << 12
		x |= inst.Dr << 9
		x |= inst.Sr1 << 6
		x |= 0b1 << 5
		x |= 0b11111

		m.write(uint16(x), inst)
	case spec.OP_TRAP:
		var x int
		x = spec.OP_TRAP << 12
		x |= int(inst.Trapvect8) & 0b11111111

		m.write(uint16(x), inst)
	default:
		m.errors.Add(inst, fmt.Sprintf("unrecognized opcode: 0b%b", inst.Opcode))
	}
}

func (m *emitter) emitFillDirective(d *ast.FillDirective) {
	m.write(uint16(d.Value), d)
}

func (m *emitter) write(x uint16, l syntax.HasLocation) {
	err := binary.Write(m.buf, binary.BigEndian, x)
	if err != nil {
		m.errors.Add(l, err.Error())
	}
}

func (m *emitter) labelToOffset(label string, maxValueMask uint16, pc uint16, loc syntax.HasLocation) int {
	if len(label) == 0 {
		m.errors.Add(loc, "label name is empty")
		return 0
	}

	labelIndex := m.tab.Lookup(label)
	offset := labelIndex - pc - 1

	if offset > maxValueMask {
		m.errors.Add(loc, "label is too far from the current instruction to fit in bit length: "+label)
	}
	return int(offset & maxValueMask)
}

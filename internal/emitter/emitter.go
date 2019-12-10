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
	m := &emitter{errors: errorList, buf: &buf}

	for _, s := range p.Statements {
		switch v := s.(type) {
		case *ast.Instruction:
			m.emitInstruction(v)
		default:
			return nil, fmt.Errorf("Unexpected statement type: %v", v)
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
}

func (m *emitter) emitInstruction(inst *ast.Instruction) {
	switch inst.Opcode {
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
			x |= inst.Imm5
		default:
			m.errors.Add(inst.Loc(), "unknown mode")
		}

		err := binary.Write(m.buf, binary.BigEndian, uint16(x))
		if err != nil {
			m.errors.Add(inst.Loc(), err.Error())
		}
	default:
		m.errors.Add(inst.Loc(), "unrecognized opcode")
	}
}

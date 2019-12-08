package emitter

import "github.com/onlyafly/oakblue/internal/ast"

// Emit emits an assembled binary image
func Emit(p *ast.Program) ([]byte, error) {
	return []byte("test byte array"), nil
}

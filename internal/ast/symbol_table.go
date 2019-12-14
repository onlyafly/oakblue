package ast

import "fmt"

type SymbolTable struct {
	symbols map[string]uint16
}

func NewSymbolTable() *SymbolTable {
	xs := make(map[string]uint16)
	return &SymbolTable{symbols: xs}
}

func (t *SymbolTable) Insert(key string, val uint16) error {
	if _, ok := t.symbols[key]; ok {
		return fmt.Errorf("attempted to insert on duplicate key")
	}

	t.symbols[key] = val
	return nil
}

func (t *SymbolTable) Lookup(key string) uint16 {
	return t.symbols[key]
}

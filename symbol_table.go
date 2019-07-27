package main

import (
	"fmt"
)

// SymbolTable maps string symbols to memory addresses
type SymbolTable struct {
	table   map[string]int
	nextRAM int
}

// AddElement inserts a new element (variable or instruction label) into
// the table and allocates a memory location to it. If -1 is passed in as the
// memory location, the next available RAM address will be allocated
func (st *SymbolTable) AddElement(sym string, addr int) {
	if addr == -1 {
		addr = st.nextRAM
		st.nextRAM++
	}
	st.table[sym] = addr
}

// Contains returns whether or not a given symbol exists in the table
func (st *SymbolTable) Contains(sym string) bool {
	if _, ok := st.table[sym]; ok {
		return true
	}
	return false
}

func (st *SymbolTable) Debug() {
	fmt.Printf("%v", st.table)
}

// GetAddress returns the memory address of a given symbol, or -1
// if it does not exist in the table
func (st *SymbolTable) GetAddress(sym string) int {
	if addr, ok := st.table[sym]; ok {
		return addr
	}
	return -1
}

// InitializeSymbolTable returns a new SymbolTable pre-populated with the built-in
// symbols. Sets nextsRAM to 16 (this is where the RAM addresses allocated for variables begin)
func InitializeSymbolTable() SymbolTable {
	pre := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for i := 0; i < 16; i++ {
		k := fmt.Sprintf("R%d", i)
		pre[k] = i
	}
	return SymbolTable{pre, 16}
}

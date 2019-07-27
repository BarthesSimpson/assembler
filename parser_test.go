package main

import (
	"os"
	"testing"

	. "github.com/franela/goblin"
)

func emptySymbolTable() *SymbolTable {
	table := make(map[string]int)
	return &SymbolTable{table, 0}
}
func TestParser(t *testing.T) {
	g := Goblin(t)
	g.Describe("Basic statement parsing", func() {
		f, _ := os.Open("./test/Test.asm")
		p := NewParser(f, emptySymbolTable())

		g.It("Should recognize an A Statement", func() {
			cmd, _ := p.parseLine("@M")
			g.Assert(cmd.ctype).Equal(A)
		})
		g.It("Should recognize an assignment C Statement", func() {
			cmd, _ := p.parseLine("D=D+A")
			g.Assert(cmd.ctype).Equal(C)
		})
		g.It("Should recognize a comparison C Statement", func() {
			cmd, _ := p.parseLine("0;JMP")
			g.Assert(cmd.ctype).Equal(C)
		})
		g.It("Should recognize an L Statement", func() {
			cmd, _ := p.parseLine("(LOOP)")
			g.Assert(cmd.ctype).Equal(L)
		})
		g.It("Should recognize a Comment Statement", func() {
			cmd, _ := p.parseLine("//D=D+A")
			g.Assert(cmd.ctype).Equal(Comment)
		})

		g.Describe("C Statement parsing", func() {
			f, _ := os.Open("./test/Test.asm")
			p := NewParser(f, emptySymbolTable())
			g.It("Should parse an assignment C Statement", func() {
				cmd, _ := p.parseCInstruction("D=D+A")
				g.Assert(cmd.comp).Equal(CompDplusA)
				g.Assert(cmd.jump).Equal(JmpNull)
				g.Assert(cmd.mloc).Equal(LocD)
				g.Assert(cmd.symbol).Equal("")
			})
			g.It("Should parse a comparison C Statement", func() {
				cmd, _ := p.parseCInstruction("D;JGT")
				g.Assert(cmd.comp).Equal(CompD)
				g.Assert(cmd.jump).Equal(JGT)
				g.Assert(cmd.mloc).Equal(LocNull)
				g.Assert(cmd.symbol).Equal("")
			})
			g.It("Should return an error for a malformed assignment C Statement", func() {
				_, err := p.parseCInstruction("D=")
				g.Assert(err != nil).IsTrue()

				_, err = p.parseCInstruction("=D")
				g.Assert(err != nil).IsTrue()
			})

			g.It("Should return an error for a malformed comparison C Statement", func() {
				_, err := p.parseCInstruction("P;JGT")
				g.Assert(err != nil).IsTrue()

				_, err = p.parseCInstruction("D;JJJ")
				g.Assert(err != nil).IsTrue()
			})
		})
		g.Describe("C Statement decomposition", func() {
			g.It("Should correctly decompose an assignment C Statement", func() {
				f, _ := os.Open("./test/CInstructions.asm")
				p := NewParser(f, emptySymbolTable())
				p.Advance()
				dest, _ := p.Dest()
				g.Assert(dest).Equal(LocD)
				comp, _ := p.Comp()
				g.Assert(comp).Equal(CompA)
				jmp, _ := p.Jump()
				g.Assert(jmp).Equal(JmpNull)
				sym, _ := p.Symbol()
				g.Assert(sym).Equal("")
			})
			g.It("Should correctly decompose a comparison C Statement", func() {
				f, _ := os.Open("./test/CInstructions.asm")
				p := NewParser(f, emptySymbolTable())
				p.Advance()
				p.Advance()
				dest, _ := p.Dest()
				g.Assert(dest).Equal(LocNull)
				comp, _ := p.Comp()
				g.Assert(comp).Equal(CompA)
				jmp, _ := p.Jump()
				g.Assert(jmp).Equal(JMP)
				sym, _ := p.Symbol()
				g.Assert(sym).Equal("")
			})
		})
		g.Describe("A Statement parsing", func() {
			f, _ := os.Open("./test/Test.asm")
			p := NewParser(f, emptySymbolTable())
			g.It("Should parse an A Statement", func() {
				cmd, _ := p.parseLine("@12345")
				g.Assert(cmd.comp).Equal(Comp0)
				g.Assert(cmd.jump).Equal(JmpNull)
				g.Assert(cmd.mloc).Equal(LocNull)
				g.Assert(cmd.symbol).Equal("12345")
			})
		})
	})

	g.Describe("Label statement parsing", func() {
		f, _ := os.Open("./test/Test.asm")
		st := InitializeSymbolTable()
		p := NewParser(f, &st)

		g.It("Should parse an L Statement", func() {
			cmd, _ := p.parseLine("(LOOP)")
			g.Assert(cmd.comp).Equal(Comp0)
			g.Assert(cmd.jump).Equal(JmpNull)
			g.Assert(cmd.mloc).Equal(LocNull)
			g.Assert(cmd.symbol).Equal("LOOP")
		})
		g.It("Should parse an A Statement with a preset label", func() {
			cmd, _ := p.parseLine("@LCL")
			g.Assert(cmd.comp).Equal(Comp0)
			g.Assert(cmd.jump).Equal(JmpNull)
			g.Assert(cmd.mloc).Equal(LocNull)
			g.Assert(cmd.symbol).Equal("1")
		})
		g.It("Should parse an A Statement with a user defined label", func() {
			p.st.AddElement("LOOP", 4)
			cmd, _ := p.parseLine("@LOOP")
			g.Assert(cmd.comp).Equal(Comp0)
			g.Assert(cmd.jump).Equal(JmpNull)
			g.Assert(cmd.mloc).Equal(LocNull)
			g.Assert(cmd.symbol).Equal("4")
		})
		g.It("Should parse an A Statement with a variable", func() {
			cmd, _ := p.parseLine("@VAR_1")
			g.Assert(cmd.comp).Equal(Comp0)
			g.Assert(cmd.jump).Equal(JmpNull)
			g.Assert(cmd.mloc).Equal(LocNull)
			g.Assert(cmd.symbol).Equal("16")
		})
	})
}

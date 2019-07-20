package main

import (
	"os"
	"testing"

	. "github.com/franela/goblin"
)

func TestParser(t *testing.T) {
	g := Goblin(t)
	g.Describe("Numbers", func() {
		f, _ := os.Open("./test/Test.asm")
		p := NewParser(f)

		g.It("Should parse an A Statement", func() {
			cmd := p.parseLine("@M")
			g.Assert(cmd.ctype).Equal(A)
		})
		g.It("Should parse a C Statement", func() {
			cmd := p.parseLine("D=D+A")
			g.Assert(cmd.ctype).Equal(C)
		})
		g.Xit("Should parse an L Statement", func() {
			g.Assert(1 + 1).Equal(3)
		})
		g.It("Should parse a Comment Statement", func() {
			cmd := p.parseLine("//D=D+A")
			g.Assert(cmd.ctype).Equal(Comment)
		})
	})
}

package main

import (
	"bufio"
	"os"
	"testing"

	. "github.com/franela/goblin"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	g := Goblin(t)
	g.Describe("Integration tests without symbols", func() {
		g.It("Adds two numbers", func() {
			CompareFiles("test/Add.asm", "test/Add.hack", "test/AddExpected.hack", t)
		})
		g.It("Finds the max of 2 numbers", func() {
			CompareFiles("test/MaxL.asm", "test/MaxL.hack", "test/MaxLExpected.hack", t)
		})
		g.It("Draws a rectangle on the screen", func() {
			CompareFiles("test/RectL.asm", "test/RectL.hack", "test/RectLExpected.hack", t)
		})
		g.It("Plays pong", func() {
			CompareFiles("test/PongL.asm", "test/PongL.hack", "test/PongLExpected.hack", t)
		})
	})
	g.Describe("Integration tests with symbols", func() {
		g.It("Finds the max of 2 numbers", func() {
			CompareFiles("test/Max.asm", "test/Max.hack", "test/MaxExpected.hack", t)
		})
		g.It("Draws a rectangle on the screen", func() {
			CompareFiles("test/Rect.asm", "test/Rect.hack", "test/RectExpected.hack", t)
		})
		g.It("Plays pong", func() {
			CompareFiles("test/Pong.asm", "test/Pong.hack", "test/PongExpected.hack", t)
		})
	})
}

func CompareFiles(infile string, outfile string, expected string, t *testing.T) {
	asm := Assembler{infile, outfile, Code{}, InitializeSymbolTable(), nil}
	asm.Convert()

	out, err := os.Open(outfile)
	if err != nil {
		t.Errorf("Unable to open output file: %s", err)
		t.FailNow()
	}
	defer out.Close()

	exp, err := os.Open(expected)
	if err != nil {
		t.Errorf("Unable to open expected file: %s", err)
		t.FailNow()
	}
	defer exp.Close()

	outscan := bufio.NewScanner(out)
	expscan := bufio.NewScanner(exp)
	l := 1
	for expscan.Scan() {
		outscan.Scan()
		if expscan.Text() != outscan.Text() {
			t.Errorf("Mismatch at line %d:\nexpected: %s\nreceived: %s", l, expscan.Text(), outscan.Text())
			t.FailNow()
		}
		l++
	}
	if outscan.Scan() {
		t.Errorf("Compiled output has extra lines. Expected %d", l)
		t.FailNow()
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	bit "github.com/golang-collections/go-datastructures/bitarray"
)

// Assembler is the main object that converts a .asm file to a binary .hack file
type Assembler struct {
	inpath  string
	outpath string
	encoder Code
	w       *bufio.Writer
}

// Convert is the main routine that processes the input file into the output file
func (asm *Assembler) Convert() {
	infile, err := os.Open(asm.inpath)
	if err != nil {
		log.Fatalf("Unable to open input file: %s", err)
	}
	defer infile.Close()

	dest, err := os.Create(asm.outpath)
	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}
	defer dest.Close()

	p := NewParser(infile)
	asm.w = bufio.NewWriter(dest)

	for p.HasMoreCommands() {
		p.Advance()
		ctype, err := p.CommandType()
		if err != nil {
			log.Fatalf("Unable to parse line %d: %s", 1, err) //TODO: get the actual line number
		}
		if ctype.IsPrintable() {
			asm.processCommand(p)
		}
	}
}

func (asm *Assembler) processCommand(p Parser) {
	cmd, err := p.CurrentCommand()
	if err != nil {
		log.Fatalf("Unable to parse line %d: %s", 1, err) //TODO: get the actual line number
	}
	if cmd.ctype == A {
		asm.writeACommand(p)
	}
	if cmd.ctype == C {
		asm.writeCCommand(p)
	}
	fmt.Println(cmd)
	return
}

func (asm *Assembler) writeACommand(p Parser) {
	sym, err := p.Symbol()
	if err != nil {
		log.Fatalf("Unable to parse line %d: %s", 1, err) //TODO: get the actual line number
	}
	ins, err := strconv.ParseInt(sym, 10, 16)
	if err != nil {
		log.Fatalf("Invalid symbol or decimal constant on line %d: %s", 1, err) //TODO: get the actual line number
	}
	asm.w.WriteString(strconv.FormatInt(ins, 2))
}

func (asm *Assembler) writeCCommand(p Parser) {

	comp, err := p.Comp()
	if err != nil {
		log.Fatalf("Unable to get Comp for C command in line %d: %s", 1, err) //TODO: get the actual line number
	}

	dest, err := p.Dest()
	if err != nil {
		log.Fatalf("Unable to get Dest for C command in line %d: %s", 1, err) //TODO: get the actual line number
	}

	jmp, err := p.Jump()
	if err != nil {
		log.Fatalf("Unable to get Dest for C command in line %d: %s", 1, err) //TODO: get the actual line number
	}

	output := bit.NewBitArray(16)
	// Fill the first three slots with 1s
	for i := uint64(0); i < 3; i++ {
		output.SetBit(i)
	}

	compBin := asm.encoder.Comp(comp)
	for i := uint64(0); i < 7; i++ {
		b, err := compBin.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", 1, err) //TODO: get the actual line number
		}
		if b == true {
			output.SetBit(i + 3)
		}
	}

	destBin := asm.encoder.Dest(dest)
	for i := uint64(0); i < 3; i++ {
		b, err := destBin.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", 1, err) //TODO: get the actual line number
		}
		if b == true {
			output.SetBit(i + 10)
		}
	}

	jmpBin := asm.encoder.Jump(jmp)
	for i := uint64(0); i < 3; i++ {
		b, err := jmpBin.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", 1, err) //TODO: get the actual line number
		}
		if b == true {
			output.SetBit(i + 13)
		}
	}

	strArr := []string{"1", "1", "1", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
	for i := uint64(0); i < 16; i++ {
		b, err := output.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", 1, err) //TODO: get the actual line number
		}
		if b {
			strArr[i] = "1"
		}
	}

	asm.w.WriteString(strings.Join(strArr, ""))
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	bit "github.com/golang-collections/go-datastructures/bitarray"
	log "github.com/sirupsen/logrus"
)

// Assembler is the main object that converts a .asm file to a binary .hack file
type Assembler struct {
	inpath  string
	outpath string
	encoder Code
	st      SymbolTable
	w       *bufio.Writer
}

// Convert is the main routine that processes the input file into the output file
func (asm *Assembler) Convert() {
	infile, err := os.Open(asm.inpath)
	if err != nil {
		log.Fatalf("Unable to open input file: %s", err)
	}

	dest, err := os.Create(asm.outpath)
	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}
	defer dest.Close()

	asm.w = bufio.NewWriter(dest)
	asm.buildSymbolTable(infile)
	infile.Close() // Do not defer this, as doing so will block translateInstructions

	infile, err = os.Open(asm.inpath)
	if err != nil {
		log.Fatalf("Unable to open input file: %s", err)
	}
	asm.translateInstructions(infile)
	infile.Close()
}

// Perform a first pass of the input file, constructing the symbol table
// that will be used in translating the assembly code into binary.
// For each label that is encountered, store the label in the table;
// for each A or C instruction, increment the RAM address that is used
// to store the next label.
func (asm *Assembler) buildSymbolTable(infile *os.File) {
	p := NewParser(infile, &asm.st)
	l := 1
	addr := 0
	for {
		p.Advance(true)
		if !p.HasMoreCommands() {
			break
		}
		ctype := p.CommandType()
		if ctype == L {
			sym, err := p.Symbol()
			if err != nil {
				log.Fatalf("Unable to retrieve symbol for line %d: %s", l, err)
			}
			asm.st.AddElement(sym, addr)
		} else if ctype == C || ctype == A {
			addr++
		}
		l++
	}
}

// Perform a second pass of the input file, during which the actual
// conversion to binary and writing of the output is performed
func (asm *Assembler) translateInstructions(infile *os.File) {
	p := NewParser(infile, &asm.st)
	l := 1
	for {
		p.Advance(false)
		if !p.HasMoreCommands() {
			break
		}
		ctype := p.CommandType()
		if ctype.IsPrintable() {
			asm.processCommand(p, l)
		}
		l++
	}
	asm.w.Flush()
}

func (asm *Assembler) processCommand(p Parser, l int) {
	cmd := p.CurrentCommand()
	if cmd.ctype == A {
		asm.writeACommand(p, l)
		return
	}
	if cmd.ctype == C {
		asm.writeCCommand(p, l)
		return
	}
}

func (asm *Assembler) writeACommand(p Parser, l int) {
	sym, err := p.Symbol()
	if err != nil {
		log.Fatalf("Unable to parse line %d: %s", l, err)
	}
	ins, err := strconv.ParseInt(sym, 10, 16)
	if err != nil {
		log.Fatalf("Invalid symbol or decimal constant on line %d: %s", l, err)
	}
	str := fmt.Sprintf("%016b\n", ins)
	log.Debug(str)
	_, err = asm.w.WriteString(str)
	if err != nil {
		log.Fatalf("Unable to write line %d: %s", l, err)
	}
}

func (asm *Assembler) writeCCommand(p Parser, l int) {

	comp, err := p.Comp()
	if err != nil {
		log.Fatalf("Unable to get Comp for C command in line %d: %s", l, err)
	}

	dest, err := p.Dest()
	if err != nil {
		log.Fatalf("Unable to get Dest for C command in line %d: %s", l, err)
	}

	jmp, err := p.Jump()
	if err != nil {
		log.Fatalf("Unable to get Dest for C command in line %d: %s", l, err)
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
			log.Fatalf("Unable to write binary output for line %d: %s", l, err)
		}
		if b == true {
			output.SetBit(i + 3)
		}
	}

	destBin := asm.encoder.Dest(dest)
	for i := uint64(0); i < 3; i++ {
		b, err := destBin.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", l, err)
		}
		if b == true {
			output.SetBit(i + 10)
		}
	}

	jmpBin := asm.encoder.Jump(jmp)
	for i := uint64(0); i < 3; i++ {
		b, err := jmpBin.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", l, err)
		}
		if b == true {
			output.SetBit(i + 13)
		}
	}

	strArr := []string{"1", "1", "1", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
	for i := uint64(0); i < 16; i++ {
		b, err := output.GetBit(i)
		if err != nil {
			log.Fatalf("Unable to write binary output for line %d: %s", l, err)
		}
		if b {
			strArr[i] = "1"
		}
	}
	log.Debug(strings.Join(strArr, ""))

	out := fmt.Sprintf("%s\n", strings.Join(strArr, ""))
	_, err = asm.w.WriteString(out)
	if err != nil {
		log.Fatalf("Unable to write line %d: %s", l, err)
	}
}

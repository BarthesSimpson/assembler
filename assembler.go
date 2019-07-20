package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Assembler is the main object that converts a .asm file to a binary .hack file
type Assembler struct {
	inpath  string
	outpath string
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
			log.Fatalf("Unable to parse line %s: %s", 1, err) //TODO: get the actual line number
		}
		if ctype.IsPrintable() {
			cmd, err := p.CurrentCommand()
			if err != nil {
				log.Fatalf("Unable to parse line %s: %s", 1, err) //TODO: get the actual line number
			}
			asm.processCommand(cmd)
		}
	}

	// scanner := bufio.NewScanner(source)
	// for scanner.Scan() {
	// 	asm.ProcessLine(scanner.Text())
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }
}

func (asm *Assembler) processCommand(cmd Command) {
	fmt.Println(cmd)
	return
}

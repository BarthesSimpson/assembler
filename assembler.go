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
func (ass *Assembler) Convert() {
	source, err := os.Open(ass.inpath)
	if err != nil {
		log.Fatalf("Unable to open input file: %s", err)
	}
	defer source.Close()

	dest, err := os.Create(ass.outpath)
	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}
	defer dest.Close()

	ass.w = bufio.NewWriter(dest)

	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		ass.ProcessLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// ProcessLine detects the type of the line and forwards it to the appropriate handler
func (ass *Assembler) ProcessLine(line string) {
	fmt.Println(line)
}

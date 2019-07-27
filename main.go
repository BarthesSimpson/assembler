package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: assemble <filepath>")
	}

	inpath := os.Args[1]
	fname := strings.Split(inpath, ".")[0]
	outpath := fmt.Sprintf("%s.hack", fname)
	asm := Assembler{inpath, outpath, Code{}, InitializeSymbolTable(), nil}
	asm.Convert()
}

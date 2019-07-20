package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: assemble <filepath>")
	}

	inpath := os.Args[1]
	fname := strings.Split(inpath, ".")[0]
	outpath := fmt.Sprintf("%s.hack", fname)

	ass := Assembler{inpath, outpath, nil}
	ass.Convert()
}

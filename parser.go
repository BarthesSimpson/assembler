package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// Command represents a single assembly command
type Command struct {
	ctype  CommandType
	comp   CompMnemonic
	jump   JumpMnemonic
	mloc   MemoryLocation
	symbol string
}

// Parser is the main object that processes the input file line by line
type Parser struct {
	infile          *os.File
	scanner         *bufio.Scanner
	currentCommand  Command
	hasMoreCommands bool
}

// NewParser is a factory that creates a parser instance for the given file
func NewParser(infile *os.File) Parser {
	scanner := bufio.NewScanner(infile)
	return Parser{infile, scanner, Command{}, true}
}

// HasMoreCommands indicates whether the entire input file has been processed
func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

// CurrentCommand returns the parsed command for the current line if it exists
func (p *Parser) CurrentCommand() (Command, error) {
	if p.currentCommand.ctype == 0 {
		return Command{}, errors.New("No command has been parsed")
	}
	return p.currentCommand, nil
}

// CommandType returns the type of the current command if it exists
func (p *Parser) CommandType() (CommandType, error) {
	if p.currentCommand.ctype == 0 {
		return 0, errors.New("No command has been parsed")
	}
	return p.currentCommand.ctype, nil
}

// Advance moves one line forward in the input file
func (p *Parser) Advance() {
	p.hasMoreCommands = p.scanner.Scan()
	if !p.hasMoreCommands {
		return
	}
	line := p.scanner.Text()
	p.hasMoreCommands = true
	cmd, err := p.parseLine(line)
	if err != nil {
		log.Fatalf("Unable to parse line %s", line)
	}
	p.currentCommand = cmd
}

func (p *Parser) parseLine(line string) (Command, error) {
	if strings.HasPrefix(line, CommentToken) {
		return Command{Comment, Comp0, JmpNull, LocNull, line[2:]}, nil
	}
	if strings.HasPrefix(line, ACmdToken) {
		return Command{A, Comp0, JmpNull, LocNull, line[1:]}, nil
	}
	if strings.ContainsAny(line, "=;") {
		cmd, err := p.parseCInstruction(line)
		if err != nil {
			return Command{}, err
		}
		return cmd, nil
	}
	return Command{L, Comp0, JmpNull, LocNull, line}, nil
}

func (p *Parser) parseCInstruction(line string) (Command, error) {
	if strings.Contains(line, "=") {
		chars := filterEmpty(strings.Split(line, "="))
		if len(chars) < 2 {
			return Command{}, fmt.Errorf("%s is not a valid C command", line)
		}
		mloc := EnumValFromString(MemoryLocationStrings, chars[0])
		if mloc == -1 {
			return Command{}, fmt.Errorf("%s is not a valid memory location", chars[0])
		}
		return Command{C, Comp0, JmpNull, MemoryLocation(mloc), chars[1]}, nil
	}

	if strings.Contains(line, ";") {
		chars := filterEmpty(strings.Split(line, ";"))
		if len(chars) < 2 {
			return Command{}, fmt.Errorf("%s is not a valid C command", line)
		}
		cmp := EnumValFromString(CompStrings, chars[0])
		if cmp == -1 {
			return Command{}, fmt.Errorf("%s is not a valid comp value", chars[0])
		}
		jmp := EnumValFromString(JumpStrings, chars[1])
		if jmp == -1 {
			return Command{}, fmt.Errorf("%s is not a valid jump expression", chars[1])
		}
		return Command{C, CompMnemonic(cmp), JumpMnemonic(jmp), LocNull, ""}, nil
	}
	return Command{}, fmt.Errorf("%s is not a C instruction", line)
}

// Symbol retrieves the symbol (variable name or constant) associated with the current command
// if it exists
func (p *Parser) Symbol() (string, error) {
	ctype, _ := p.CommandType()
	if !ctype.IsPrintable() {
		return "", errors.New("only A commands and C commands can contain symbols")
	}
	cmd, err := p.CurrentCommand()
	if err != nil {
		return "", errors.New("the parser has no command loaded")
	}
	return cmd.symbol, nil
}

// Dest retrieves the memory location where the current C command should write its output
func (p *Parser) Dest() (MemoryLocation, error) {
	cmd, err := p.CurrentCommand()
	if err != nil || cmd.ctype != C {
		return LocNull, errors.New("the parser has no C command loaded")
	}
	return cmd.mloc, nil
}

// Comp retrieves the comp mnemonic of the current C command
func (p *Parser) Comp() (CompMnemonic, error) {
	cmd, err := p.CurrentCommand()
	if err != nil || cmd.ctype != C {
		return Comp0, errors.New("the parser has no C command loaded")
	}
	return cmd.comp, nil
}

// Jump retrieves the jump mnemonic of the current C command
func (p *Parser) Jump() (JumpMnemonic, error) {
	cmd, err := p.CurrentCommand()
	if err != nil || cmd.ctype != C {
		return JmpNull, errors.New("the parser has no C command loaded")
	}
	return cmd.jump, nil
}

func filterEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

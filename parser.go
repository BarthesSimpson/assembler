package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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
	st              *SymbolTable
	scanner         *bufio.Scanner
	currentCommand  Command
	hasMoreCommands bool
}

// NewParser is a factory that creates a parser instance for the given file
func NewParser(infile *os.File, st *SymbolTable) Parser {
	scanner := bufio.NewScanner(infile)
	return Parser{infile, st, scanner, Command{}, true}
}

// HasMoreCommands indicates whether the entire input file has been processed
func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

// CurrentCommand returns the parsed command for the current line if it exists
func (p *Parser) CurrentCommand() Command {
	return p.currentCommand
}

// CommandType returns the type of the current command if it exists
func (p *Parser) CommandType() CommandType {
	return p.currentCommand.ctype
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
	line = stripInlineComments(line)
	if strings.HasPrefix(line, LabelToken) {
		return Command{L, Comp0, JmpNull, LocNull, line[1 : len(line)-1]}, nil
	}
	if strings.HasPrefix(line, ACmdToken) {
		cmd, err := p.parseAInstruction(line)
		if err != nil {
			return Command{}, err
		}
		return cmd, nil
	}
	if strings.ContainsAny(line, "=;") {
		cmd, err := p.parseCInstruction(line)
		if err != nil {
			return Command{}, err
		}
		return cmd, nil
	}
	return Command{CmdNull, Comp0, JmpNull, LocNull, ""}, nil
}

func (p *Parser) parseAInstruction(line string) (Command, error) {
	sym := line[1:]
	// If symbol is an integer literal, we can just return is as is
	if _, err := strconv.Atoi(sym); err == nil {
		return Command{A, Comp0, JmpNull, LocNull, sym}, nil
	}
	// If symbol is already in the symbol table, just resolve it;
	// Otherwise, insert it at the next available RAM location
	if !p.st.Contains(sym) {
		p.st.AddElement(sym, -1)
	}
	addr := fmt.Sprintf("%d", p.st.GetAddress(sym))
	return Command{A, Comp0, JmpNull, LocNull, addr}, nil
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
		comp := EnumValFromString(CompStrings, chars[1])
		if comp == -1 {
			return Command{}, fmt.Errorf("%s is not a valid memory location", chars[0])
		}
		return Command{C, CompMnemonic(comp), JmpNull, MemoryLocation(mloc), ""}, nil
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
	ctype := p.CommandType()
	if !ctype.IsPrintable() && ctype != L {
		return "", errors.New("only A commands and C commands can contain symbols")
	}
	cmd := p.CurrentCommand()
	return cmd.symbol, nil
}

// Dest retrieves the memory location where the current C command should write its output
func (p *Parser) Dest() (MemoryLocation, error) {
	cmd := p.CurrentCommand()
	if cmd.ctype != C {
		return LocNull, errors.New("the parser has no C command loaded")
	}
	return cmd.mloc, nil
}

// Comp retrieves the comp mnemonic of the current C command
func (p *Parser) Comp() (CompMnemonic, error) {
	cmd := p.CurrentCommand()
	if cmd.ctype != C {
		return Comp0, errors.New("the parser has no C command loaded")
	}
	return cmd.comp, nil
}

// Jump retrieves the jump mnemonic of the current C command
func (p *Parser) Jump() (JumpMnemonic, error) {
	cmd := p.CurrentCommand()
	if cmd.ctype != C {
		return JmpNull, errors.New("the parser has no C command loaded")
	}
	return cmd.jump, nil
}

// Helper functions

func filterEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func stripInlineComments(line string) string {
	split := strings.Split(line, CommentToken)
	return strings.Trim(split[0], " ")
}

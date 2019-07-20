package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Command represents a single assembly command
type Command struct {
	ctype CommandType
}

// Parser is the main object that processes the input file line by line
type Parser struct {
	infile          *os.File
	scanner         *bufio.Scanner
	currentCommand  Command
	hasMoreCommands bool
}

// CommandType is an integer enum type
type CommandType int

// Enum for the possible types of command:
// A is a memory assignment
// C is an operation
// L is a symbol or variable assignment
// Comment is a commented line that will be ignored
const (
	A CommandType = iota
	C
	L
	Comment
)

// IsPrintable determines whether the command is a printable command (a or c type)
// or a non-printable (comment or pseudo-command)
func (cmd CommandType) IsPrintable() bool {
	return cmd < 2
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
	fmt.Print("CurrentCommand is not yet implemented")
	return Command{}, nil
}

// CommandType returns the type of the current command if it exists
func (p *Parser) CommandType() (CommandType, error) {
	fmt.Print("CommandType is not yet implemented")
	return 3, nil
}

// Advance moves one line forward in the input file
func (p *Parser) Advance() {
	p.hasMoreCommands = p.scanner.Scan()
	if !p.hasMoreCommands {
		return
	}
	line := p.scanner.Text()
	p.hasMoreCommands = true
	p.currentCommand = p.parseLine(line)
}

func (p *Parser) parseLine(line string) Command {
	if strings.HasPrefix(line, CommentToken) {
		return Command{Comment}
	}
	if strings.HasPrefix(line, ACmdToken) {
		return Command{A}
	}
	// TODO: distinguish between c commands and l commands
	return Command{C}
}

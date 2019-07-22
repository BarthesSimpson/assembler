package main

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

// CommandTypeStrings enables converting a CommandType to and from its string representation
var CommandTypeStrings = []string{"A", "C", "L", "Comment"}

// IsPrintable determines whether the command is a printable command (a or c type)
// or a non-printable (comment or pseudo-command)
func (cmd CommandType) IsPrintable() bool {
	return cmd < 2
}

// MemoryLocation is an integer enum type
type MemoryLocation int

// Enum for the possible memory locations that a C command
// should write its output to
const (
	LocNull MemoryLocation = iota
	LocM
	LocD
	LocMD
	LocA
	LocAM
	LocAD
	LocAMD
)

// MemoryLocationStrings enables converting a MemoryLocation to and from its string representation
var MemoryLocationStrings = []string{"null", "M", "D", "MD", "A", "AM", "AD", "AMD"}

// Comp is an integer enum type
type Comp int

// Enum for the possible comp mnemonics that a C command
// can encode
const (
	Comp0 Comp = iota
	Comp1
	CompMinus1
	CompD
	CompA
	CompNegD
	CompNegA
	CompMinusD
	CompMinusA
	CompDplus1
	CompAplus1
	CompDminus1
	CompAminus1
	CompDplusA
	CompDminusA
	CompAminusD
	CompDandA
	CompAorA
	CompM
	CompNegM
	CompAminusM
	CompMplus1
	CompMminus1
	CompDplusM
	CompDminusM
	CompMminusD
	CompDandM
	CompDorM
)

// CompStrings enables converting a Comp to and from its string representation
var CompStrings = []string{"0", "1", "-1", "D", "A", "!D", "!A", "-D", "-A", "D+1", "A+1", "D-1", "A-1", "D+A", "D-A", "A-D", "D&A", "D|A", "M", "!M", "-M", "M+1", "M-1", "D+M", "D-M", "M-D", "D&M", "D|M"}

// Jump is an integer enum type
type Jump int

// Enum for the possible jump mnemonics that a C command
// can encode
const (
	JmpNull Jump = iota
	JGT
	FEQ
	JGE
	JLT
	JNE
	JLE
	JMP
)

// JumpStrings enables converting a Jump to and from its string representation
var JumpStrings = []string{"null", "JGT", "FEQ", "JGE", "JLT", "JNE", "JLE", "JMP"}

// EnumValFromString enables converting a string into an enum value
func EnumValFromString(enumStrings []string, searchVal string) int {
	for i, s := range enumStrings {
		if s == searchVal {
			return i
		}
	}
	return -1
}

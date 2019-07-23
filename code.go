package main

import bit "github.com/golang-collections/go-datastructures/bitarray"

// Code converts machine language tokens into binary
type Code struct{}

// Dest converts a MemoryLocation mnemonic into its binary representation
func (c *Code) Dest(mloc MemoryLocation) bit.BitArray {
	return bit.NewBitArray(3)
}

// Comp converts a Comp mnemonic into its binary representation
func (c *Code) Comp(comp CompMnemonic) bit.BitArray {
	out := bit.NewBitArray(7)
	if comp == Comp0 {
		return setBits(out, []bool{false, true, false, true, false, true, false})
	}
	if comp == Comp1 {
		return setBits(out, []bool{false, true, true, true, true, true, true})
	}
	if comp == CompMinus1 {
		return setBits(out, []bool{false, true, true, true, false, true, false})
	}
	if comp == CompD {
		return setBits(out, []bool{false, false, false, true, true, false, false})
	}
	if comp == CompA {
		return setBits(out, []bool{false, true, true, false, false, false, false})
	}
	if comp == CompNegD {
		return setBits(out, []bool{false, false, false, true, true, false, true})
	}
	if comp == CompNegA {
		return setBits(out, []bool{false, true, true, false, false, false, true})
	}
	if comp == CompMinusD {
		return setBits(out, []bool{false, false, false, true, true, true, true})
	}
	if comp == CompMinusA {
		return setBits(out, []bool{false, true, true, false, false, true, true})
	}
	if comp == CompDplus1 {
		return setBits(out, []bool{false, false, true, false, false, true, true})
	}
	if comp == CompAplus1 {
		return setBits(out, []bool{false, true, true, false, true, true, true})
	}
	if comp == CompDminus1 {
		return setBits(out, []bool{false, true, true, false, false, true, false})
	}
	if comp == CompAminus1 {
		return setBits(out, []bool{false, true, true, false, true, true, true})
	}
	if comp == CompDplusA {
		return setBits(out, []bool{false, false, false, false, false, true, false})
	}
	if comp == CompDminusA {
		return setBits(out, []bool{false, false, true, false, false, true, true})
	}
	if comp == CompAminusD {
		return setBits(out, []bool{false, false, false, false, true, true, true})
	}
	if comp == CompDandA {
		return setBits(out, []bool{false, false, false, false, false, false, false})
	}
	if comp == CompDorA {
		return setBits(out, []bool{false, false, true, false, true, false, true})
	}
	if comp == CompM {
		return setBits(out, []bool{true, true, true, false, false, false, false})
	}
	if comp == CompNegM {
		return setBits(out, []bool{true, true, true, false, false, false, true})
	}
	if comp == CompMinusM {
		return setBits(out, []bool{true, true, true, false, false, true, true})
	}
	if comp == CompMplus1 {
		return setBits(out, []bool{true, true, true, false, true, true, true})
	}
	if comp == CompMminus1 {
		return setBits(out, []bool{true, true, true, false, true, true, true})
	}
	if comp == CompDplusM {
		return setBits(out, []bool{true, false, false, false, false, true, false})
	}
	if comp == CompDminusM {
		return setBits(out, []bool{true, false, true, false, false, true, true})
	}
	if comp == CompMminusD {
		return setBits(out, []bool{true, false, false, false, true, true, true})
	}
	if comp == CompDandM {
		return setBits(out, []bool{true, false, false, false, false, false, false})
	}
	if comp == CompDorM {
		return setBits(out, []bool{true, false, true, false, true, false, true})
	}
	return out
}

// Jump converts a Jump mnemonic into its binary representation
func (c *Code) Jump(jmp JumpMnemonic) bit.BitArray {
	return bit.NewBitArray(3)
}

func setBits(arr bit.BitArray, bits []bool) bit.BitArray {
	for i, b := range bits {
		if b == true {
			arr.SetBit(uint64(i))
		}
	}
	return arr
}

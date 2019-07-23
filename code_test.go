package main

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestCode(t *testing.T) {
	g := Goblin(t)
	g.Describe("Comp translation", func() {
		c := Code{}
		g.It("Generates unique binary outputs for every Comp mnemonic", func() {
			m := map[string]string{}
			for i, s := range CompStrings {
				cmp := CompMnemonic(i)
				bits := fmt.Sprintf("%v", c.Comp(cmp))
				if prev, exists := m[bits]; exists == true {
					err := fmt.Sprintf("Same bit value found for %s and %s", s, prev)
					g.Assert(false).IsTrue(err)
				}
				m[bits] = s
			}
		})
	})

	g.Describe("Dest translation", func() {
		c := Code{}
		g.It("Generates unique binary outputs for every Dest mnemonic", func() {
			m := map[string]string{}
			for i, s := range MemoryLocationStrings {
				mloc := MemoryLocation(i)
				bits := fmt.Sprintf("%v", c.Dest(mloc))
				if prev, exists := m[bits]; exists == true {
					err := fmt.Sprintf("Same bit value found for %s and %s", s, prev)
					g.Assert(false).IsTrue(err)
				}
				m[bits] = s
			}
		})
	})
}

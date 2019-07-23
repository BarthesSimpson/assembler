package main

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
	bit "github.com/golang-collections/go-datastructures/bitarray"
)

func TestCode(t *testing.T) {
	g := Goblin(t)
	g.Describe("Comp translation", func() {
		c := Code{}
		g.It("Generates unique binary outputs for every Comp mnemonic", func() {
			m := map[bit.BitArray]string{}
			for i, s := range CompStrings {
				cmp := CompMnemonic(i)
				bits := c.Comp(cmp)
				if prev, exists := m[bits]; exists == true {
					err := fmt.Sprintf("Same bit value found for %s and %s", s, prev)
					g.Assert(false).IsTrue(err)
				}
				m[bits] = s
			}
		})
	})
}

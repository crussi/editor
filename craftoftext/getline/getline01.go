package main

// getline01.go
//
// GetLine Version One — from "The Craft of Text Editing" by Craig Finseth.
//
// Accepts printable characters and the Enter key only.
// All other keys cause a beep. There is no way to correct mistakes.

import (
	"fmt"
	"unicode"
)

type getlineV1 struct{}

func (g getlineV1) GetLine(prompt string, buffer []byte) bool {
	if len(buffer) < 2 {
		return false // safety check
	}

	fmt.Printf("%s: ", prompt)

	pos := 0
	for {
		key := keyGet()

		if unicode.IsPrint(rune(key)) {
			if pos >= len(buffer)-1 {
				beep() // buffer full
			} else {
				buffer[pos] = byte(key)
				pos++
				fmt.Printf("%c", key)
			}
		} else if key == keyEnter {
			buffer[pos] = 0 // NUL-terminate
			fmt.Println()
			return true
		} else {
			beep() // unknown key
		}
	}
}

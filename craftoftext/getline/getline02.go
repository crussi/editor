package main

// getline02.go
//
// GetLine Version Two — from "The Craft of Text Editing" by Craig Finseth.
//
// Adds Backspace/Delete editing over Version One.
// The user can now erase the last character typed.

import (
	"fmt"
	"unicode"
)

type getlineV2 struct{}

func (g getlineV2) GetLine(prompt string, buffer []byte) bool {
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
		} else {
			switch key {
			case keyBack:
				if pos > 0 {
					pos--
					fmt.Print("\b \b") // move back, overwrite with space, move back again
				}
				// silently ignore backspace at start of line

			case keyEnter:
				buffer[pos] = 0 // NUL-terminate
				fmt.Println()
				return true

			default:
				beep() // unknown key
			}
		}
	}
}

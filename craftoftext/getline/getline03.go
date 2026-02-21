package main

// getline03.go
//
// GetLine Version Three — from "The Craft of Text Editing" by Craig Finseth.
//
// Adds support for a pre-loaded default value.
// If the user presses Enter immediately, the existing buffer contents
// are returned unchanged.  The first printable key or Backspace clears
// the default and starts fresh.

import (
	"fmt"
	"unicode"
)

type getlineV3 struct{}

func (g getlineV3) GetLine(prompt string, buffer []byte) bool {
	if len(buffer) < 2 {
		return false // safety check
	}

	// wasKey: has the user started typing yet?
	// Until they do, the default (already in buffer) is shown but will
	// be wiped on the first keystroke.
	wasKey := false

	for {
		// Redisplay the whole prompt + current buffer on every iteration.
		fmt.Print("\r\033[2K") // carriage return, erase line (ANSI)
		fmt.Printf("%s: %s", prompt, cstring(buffer))

		key := keyGet()

		if unicode.IsPrint(rune(key)) {
			if !wasKey {
				buffer[0] = 0 // clear the default
				wasKey = true
			}
			pos := clen(buffer)
			if pos >= len(buffer)-1 {
				beep()
			} else {
				buffer[pos] = byte(key)
				buffer[pos+1] = 0
			}

		} else {
			switch key {
			case keyBack:
				if !wasKey {
					buffer[0] = 0 // clear the default
					wasKey = true
				}
				pos := clen(buffer)
				if pos > 0 {
					buffer[pos-1] = 0
				}

			case keyEnter:
				fmt.Println()
				return true

			default:
				beep()
			}
		}
	}
}

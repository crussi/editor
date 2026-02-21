package main

// getline05.go
//
// GetLine Version Five — from "The Craft of Text Editing" by Craig Finseth.
// Ch 1 Question 1 - Modify the latest version of Get_Line to accept only numeric responses. What sort of error messages should be given? (Easy)
//

import (
	"fmt"
	"unicode"
)

type getlineV6 struct {
	Allowed []string
}

// The function signature is unchanged from previous versions so that
// main.go needs no edits.  The default value is whatever is already
// in buffer when the function is called; pre-load it before calling
// if you want a default.
func (g getlineV6) GetLine(prompt string, buffer []byte) bool {
	if len(buffer) < 2 {
		return false // safety check
	}

	// Save original buffer contents so Ctrl-R can restore the default.
	saved := make([]byte, len(buffer))
	copy(saved, buffer)

	// cursor — index of the character the cursor sits on (0-based).
	// wasKey — has the user pressed anything yet?
	// insert — true = insert mode, false = replace mode.
	cursor := clen(buffer) // start at end of any pre-loaded default
	wasKey := false
	insert := true

	// Print helper bar once before entering the loop
	printHelper := func() {
		mode := "INS"
		if !insert {
			mode = "REP"
		}
		fmt.Print("\r\033[2K") // clear line
		fmt.Printf("[%s] ← → Home End | BS Del | Ctrl-U clear | Ctrl-R default | Ctrl-P quote | Ctrl-G cancel | Ctrl-L redisplay", mode)
		fmt.Println()
	}

	// Print helper bar initially
	printHelper()

	for {
		// ── Redisplay input line only ─────────────────────────────────────
		fmt.Print("\r\033[2K") // clear input line
		fmt.Printf("%s: %s", prompt, cstring(buffer))

		// Move cursor to correct column
		col := len(prompt) + 2 + cursor
		fmt.Printf("\r\033[%dC", col)

		// ── Read key ─────────────────────────────────────────────────────
		key := keyGetExt()

		// If insert/replace mode changed, redraw helper bar
		if key == keyInsToggle {
			insert = !insert
			// Move cursor up one line, redraw helper, move back down
			fmt.Print("\033[1A") // up
			printHelper()
			fmt.Print("\033[1B") // down
			continue
		}

		// ── Printable character ──────────────────────────────────────────
		if key > 0 && unicode.IsPrint(rune(key)) {

			if !wasKey {
				buffer[0] = 0
				cursor = 0
				wasKey = true
			}
			if insert {
				if clen(buffer) >= len(buffer)-1 {
					beep()
				} else {
					insertChar(buffer, cursor, byte(key))
					cursor++
				}
			} else {
				if cursor >= len(buffer)-1 {
					beep()
				} else {
					if cursor == clen(buffer) {
						buffer[cursor+1] = 0
					}
					buffer[cursor] = byte(key)
					cursor++
				}
			}
			continue
		}

		// ── Control / special keys ───────────────────────────────────────
		switch key {

		case keyBack:
			if !wasKey {
				buffer[0] = 0
				cursor = 0
				wasKey = true
			}
			if cursor > 0 {
				deleteChar(buffer, cursor-1)
				cursor--
			}

		case keyDel:
			if cursor < clen(buffer) {
				deleteChar(buffer, cursor)
			} else {
				beep()
			}

		case keyLeft:
			wasKey = true
			if cursor > 0 {
				cursor--
			}

		case keyRight:
			wasKey = true
			if cursor < clen(buffer) {
				cursor++
			}

		case keyHome:
			wasKey = true
			cursor = 0

		case keyEnd:
			wasKey = true
			cursor = clen(buffer)

		case keyEnter:
			response := cstring(buffer)

			// Validate against allowed list
			if !g.isAllowed(response) {
				beep()
				fmt.Printf("\nInvalid response: %q\n", response)
				fmt.Printf("Allowed values: %v\n", g.Allowed)

				// Clear the buffer and reset cursor
				buffer[0] = 0
				cursor = 0
				wasKey = false

				// Redisplay helper bar
				fmt.Print("\033[1A")
				printHelper()
				fmt.Print("\033[1B")

				continue
			}

			fmt.Println()
			return true

		case keyCtrlG:
			fmt.Println()
			return false

		case keyCtrlU:
			buffer[0] = 0
			cursor = 0
			wasKey = true

		case keyCtrlR:
			copy(buffer, saved)
			cursor = clen(buffer)
			wasKey = false

		case keyCtrlP:
			if !wasKey {
				buffer[0] = 0
				cursor = 0
				wasKey = true
			}
			literal := keyGetExt()
			if clen(buffer) >= len(buffer)-1 {
				beep()
			} else {
				insertChar(buffer, cursor, byte(literal))
				cursor++
			}

		case keyCtrlL:
			// Redisplay helper + input line
			fmt.Print("\033[1A") // up
			printHelper()
			fmt.Print("\033[1B") // down

		default:
			beep()
		}
	}
}

func (g getlineV6) isAllowed(s string) bool {
	for _, v := range g.Allowed {
		if s == v {
			return true
		}
	}
	return false
}

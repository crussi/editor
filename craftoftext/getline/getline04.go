package main

// getline04.go
//
// GetLine Version Four — from "The Craft of Text Editing" by Craig Finseth.
//

import (
	"fmt"
	"unicode"
)

// Additional key codes needed by Version 4.
// These are synthetic values returned by keyGetV4() when it decodes
// multi-byte escape sequences from arrow keys, Home, End, and Delete.
const (
	keyLeft      = -1 // left arrow
	keyRight     = -2 // right arrow
	keyHome      = -3 // Home key
	keyEnd       = -4 // End key
	keyDel       = -5 // Delete key (forward delete)
	keyCtrlG     = 7  // Ctrl-G  — cancel
	keyCtrlL     = 12 // Ctrl-L  — redisplay
	keyCtrlP     = 16 // Ctrl-P  — quote next character literally
	keyCtrlU     = 21 // Ctrl-U  — clear line
	keyCtrlR     = 18 // Ctrl-R  — restore default
	keyInsToggle = 26 // Ctrl-Z  — toggle insert / replace mode
)

type getlineV4 struct{}

// GetLine — Version Four.
//
// Adds over Version Three:
//   - left / right cursor movement (arrow keys)
//   - Home / End keys
//   - Forward delete (Del key)
//   - Insert / replace mode toggle (Ctrl-Z)
//   - Quote next character literally (Ctrl-P)
//   - Clear line (Ctrl-U)
//   - Restore default (Ctrl-R)
//   - Redisplay (Ctrl-L)
//   - Cancel / abort (Ctrl-G) — returns false
//
// The function signature is unchanged from previous versions so that
// main.go needs no edits.  The default value is whatever is already
// in buffer when the function is called; pre-load it before calling
// if you want a default.
func (g getlineV4) GetLine(prompt string, buffer []byte) bool {
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

	for {
		// ── Redisplay ────────────────────────────────────────────────
		fmt.Print("\r\033[2K") // carriage return + erase whole line
		fmt.Printf("%s: %s", prompt, cstring(buffer))

		// Move the terminal cursor to the right column.
		col := len(prompt) + 2 + cursor
		fmt.Printf("\r\033[%dC", col)

		// ── Read one key ─────────────────────────────────────────────
		key := keyGetV4()

		// ── Printable character ──────────────────────────────────────
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
				// Replace mode: overwrite or append.
				if cursor >= len(buffer)-1 {
					beep()
				} else {
					if cursor == clen(buffer) {
						buffer[cursor+1] = 0 // extend NUL terminator
					}
					buffer[cursor] = byte(key)
					cursor++
				}
			}
			continue
		}

		// ── Control / special keys ───────────────────────────────────
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
			fmt.Printf("\r\033[%dC", len(prompt)+2+clen(buffer))
			fmt.Println()
			return true

		case keyCtrlG:
			fmt.Println()
			return false // cancel

		case keyCtrlU:
			buffer[0] = 0
			cursor = 0
			wasKey = true

		case keyCtrlR:
			// Restore the original default value.
			copy(buffer, saved)
			cursor = clen(buffer)
			wasKey = false

		case keyCtrlP:
			// Quote: insert the very next keypress literally,
			// even if it is a control character.
			if !wasKey {
				buffer[0] = 0
				cursor = 0
				wasKey = true
			}
			literal := keyGetV4()
			if clen(buffer) >= len(buffer)-1 {
				beep()
			} else {
				insertChar(buffer, cursor, byte(literal))
				cursor++
			}

		case keyInsToggle:
			insert = !insert

		case keyCtrlL:
			// Redisplay — loop will redraw at the top automatically.

		default:
			beep()
		}
	}
}

// ── Buffer helpers ────────────────────────────────────────────────────────────

// insertChar inserts byte ch at position pos in a NUL-terminated buffer,
// shifting everything from pos onward one place to the right.
func insertChar(buffer []byte, pos int, ch byte) {
	length := clen(buffer)
	for i := length; i >= pos; i-- {
		buffer[i+1] = buffer[i]
	}
	buffer[pos] = ch
}

// deleteChar removes the byte at position pos from a NUL-terminated buffer,
// shifting everything after pos one place to the left.
func deleteChar(buffer []byte, pos int) {
	length := clen(buffer)
	for i := pos; i < length; i++ {
		buffer[i] = buffer[i+1]
	}
}

// ── Extended keyGet for Version 4 ────────────────────────────────────────────

// keyGetV4 wraps keyGet() and additionally decodes the multi-byte escape
// sequences that terminals send for arrow keys, Home, End, and Delete,
// returning the synthetic negative constants defined at the top of this file.
//
// Sequences decoded:
//
//	ESC [ A    up arrow    (mapped to Home)
//	ESC [ B    down arrow  (mapped to End)
//	ESC [ C    right arrow
//	ESC [ D    left arrow
//	ESC [ 1 ~  Home
//	ESC [ 3 ~  Delete (forward delete)
//	ESC [ 4 ~  End
//	ESC [ 7 ~  Home (rxvt)
//	ESC [ 8 ~  End  (rxvt)
//	ESC [ H    Home (alternate)
//	ESC [ F    End  (alternate)
func keyGetV4() int {
	ch := keyGet()
	if ch != 27 { // 27 = ESC
		return ch
	}

	// Got ESC — read the next byte to see if this is a CSI sequence.
	next := keyGet()
	if next != '[' {
		return 27 // bare ESC or unrecognised sequence
	}

	// Read the character(s) after ESC [
	ch2 := keyGet()
	switch ch2 {
	case 'A':
		return keyHome // up arrow → Home
	case 'B':
		return keyEnd // down arrow → End
	case 'C':
		return keyRight
	case 'D':
		return keyLeft
	case 'H':
		return keyHome
	case 'F':
		return keyEnd
	case '1':
		keyGet() // consume '~'
		return keyHome
	case '3':
		keyGet() // consume '~'
		return keyDel
	case '4':
		keyGet() // consume '~'
		return keyEnd
	case '7':
		keyGet() // consume '~'
		return keyHome
	case '8':
		keyGet() // consume '~'
		return keyEnd
	}

	beep()
	return 0
}

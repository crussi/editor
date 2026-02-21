// getline_v1.go
//
// A complete, executable Go implementation of the Get_Line Version One
// function from "The Craft of Text Editing" by Craig Finseth.
//
// Version One accepts printable characters and the Enter key only.
// All other keys cause a beep. There is no way to correct mistakes.
//
// Run with:  go run getline_v1.go

package main

import (
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
)

// ---------------------------------------------------------------
// Constants
// ---------------------------------------------------------------

const (
	keyEnter = 13 // Carriage Return
	bufSize  = 80 // max input length (including NUL terminator)
)

// ---------------------------------------------------------------
// Low-level terminal helpers
// ---------------------------------------------------------------

// keyGet reads a single raw keystroke from stdin and returns it as an int.
// It puts the terminal into raw mode for the duration of the read so that
// characters are delivered one at a time without waiting for Enter, and
// without the OS echoing them back automatically.
func keyGet() int {
	// Save current terminal state and restore it when done.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// If we can't go raw (e.g. stdin is a pipe), fall back to
		// a simple buffered read.
		var b [1]byte
		os.Stdin.Read(b[:])
		return int(b[0])
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b [1]byte
	os.Stdin.Read(b[:])
	return int(b[0])
}

// beep writes the ASCII BEL character to stdout.
func beep() {
	fmt.Print("\a")
}

// ---------------------------------------------------------------
// Get_Line — Version One
// (faithfully translated from Finseth's C original)
// ---------------------------------------------------------------

// GetLine prompts the user and reads a line of input into buffer.
//
// Parameters:
//
//	prompt  – text shown before the ": " separator
//	buffer  – byte slice that receives the input; must be at least 2 bytes
//
// Returns true when the user presses Enter (buffer holds NUL-terminated
// input), or false if buffer is too short to be useful.
//
// Behaviour matches Version One exactly:
//   - Printable characters are echoed and appended; overflow → beep.
//   - Enter (CR) terminates input and NUL-terminates the buffer.
//   - Any other key causes a beep; there is no backspace/delete.
func GetLine(prompt string, buffer []byte) bool {
	if len(buffer) < 2 {
		return false // safety check
	}

	fmt.Printf("%s: ", prompt)

	pos := 0
	for {
		key := keyGet()

		if unicode.IsPrint(rune(key)) {
			if pos >= len(buffer)-1 {
				beep()
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
			beep()
		}
	}
}

// ---------------------------------------------------------------
// main — simple demonstration
// ---------------------------------------------------------------

func main() {
	buffer := make([]byte, bufSize)

	fmt.Println("Get_Line Version One demo")
	fmt.Println("(only printable characters accepted; no backspace)")
	fmt.Println()

	ok := GetLine("Enter some text", buffer)
	if !ok {
		fmt.Fprintln(os.Stderr, "GetLine failed (buffer too small)")
		os.Exit(1)
	}

	// Find the NUL terminator to get the Go string.
	end := 0
	for end < len(buffer) && buffer[end] != 0 {
		end++
	}
	input := string(buffer[:end])

	fmt.Printf("\nYou entered: %q\n", input)
	fmt.Printf("Length     : %d characters\n", len(input))
}

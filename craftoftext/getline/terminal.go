// terminal.go
//
// Low-level terminal helpers shared by all GetLine versions.

package main

import (
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

// ─────────────────────────────────────────────────────────────
// Shared key constants (used by V4, V5, and future versions)
// ─────────────────────────────────────────────────────────────

const (
	keyEnter = 13
	keyBack  = 127
	bufSize  = 80

	keyLeft  = -1
	keyRight = -2
	keyHome  = -3
	keyEnd   = -4
	keyDel   = -5

	keyCtrlG     = 7
	keyCtrlL     = 12
	keyCtrlP     = 16
	keyCtrlU     = 21
	keyCtrlR     = 18
	keyInsToggle = 26
)

// ─────────────────────────────────────────────────────────────
// Raw key input
// ─────────────────────────────────────────────────────────────

func keyGet() int {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		var b [1]byte
		os.Stdin.Read(b[:])
		return int(b[0])
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b [1]byte
	os.Stdin.Read(b[:])
	return int(b[0])
}

// ─────────────────────────────────────────────────────────────
// Extended key decoder (arrow keys, Home, End, Delete)
// ─────────────────────────────────────────────────────────────

func keyGetExt() int {
	ch := keyGet()
	if ch != 27 { // ESC
		return ch
	}

	next := keyGet()
	if next != '[' {
		return 27
	}

	ch2 := keyGet()
	switch ch2 {
	case 'A':
		return keyHome
	case 'B':
		return keyEnd
	case 'C':
		return keyRight
	case 'D':
		return keyLeft
	case 'H':
		return keyHome
	case 'F':
		return keyEnd
	case '1':
		keyGet()
		return keyHome
	case '3':
		keyGet()
		return keyDel
	case '4':
		keyGet()
		return keyEnd
	case '7':
		keyGet()
		return keyHome
	case '8':
		keyGet()
		return keyEnd
	}

	beep()
	return 0
}

// ─────────────────────────────────────────────────────────────
// Buffer helpers (shared by V4, V5, future versions)
// ─────────────────────────────────────────────────────────────

func clen(b []byte) int {
	for i, v := range b {
		if v == 0 {
			return i
		}
	}
	return len(b)
}

func insertChar(buffer []byte, pos int, ch byte) {
	length := clen(buffer)
	for i := length; i >= pos; i-- {
		buffer[i+1] = buffer[i]
	}
	buffer[pos] = ch
}

func deleteChar(buffer []byte, pos int) {
	length := clen(buffer)
	for i := pos; i < length; i++ {
		buffer[i] = buffer[i+1]
	}
}

// ─────────────────────────────────────────────────────────────
// Beep + cstring
// ─────────────────────────────────────────────────────────────

func beep() {
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "-c", "[console]::beep(800,200)").Start()
	} else {
		exec.Command("echo", "-e", "\a").Start()
	}
}

func cstring(b []byte) string {
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

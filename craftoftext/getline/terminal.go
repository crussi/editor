// terminal.go
//
// Low-level terminal helpers shared by all GetLine versions.
// Nothing in here needs to change as you add new versions.

package main

import (
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

const (
	keyEnter = 13  // Carriage Return
	keyBack  = 127 // Backspace / Delete
	bufSize  = 80  // default input buffer size
)

// keyGet reads one raw keystroke from stdin and returns it as an int.
// Puts the terminal into raw mode so each keypress arrives immediately,
// one at a time, without the OS buffering or echoing it.
func keyGet() int {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// Not a real terminal (e.g. piped input) — plain read fallback.
		var b [1]byte
		os.Stdin.Read(b[:])
		return int(b[0])
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b [1]byte
	os.Stdin.Read(b[:])
	return int(b[0])
}

// beep sounds the terminal bell.
// func beep() {
// 	fmt.Print("\a")
// }

// very slow, but will keep for now
func beep() {
	if runtime.GOOS == "windows" {
		exec.Command("powershell", "-c", "[console]::beep(800,200)").Start()
	} else {
		exec.Command("echo", "-e", "\a").Start()
	}
}

// beep sounds the terminal bell, with OS-specific fallbacks.
// func beep() {
// 	// Try the classic terminal BEL first
// 	fmt.Print("\a")

// 	// OS-specific guaranteed beep
// 	switch runtime.GOOS {
// 	case "windows":
// 		exec.Command("powershell", "-c", "[console]::beep(800,200)").Run()

// 	case "darwin": // macOS
// 		exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Run()

// 	case "linux":
// 		exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/bell.oga").Run()
// 	}
// }

// cstring converts a NUL-terminated byte slice to a regular Go string.
func cstring(b []byte) string {
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

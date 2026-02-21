package main

import (
	"fmt"
	"os"
)

// Liner is the interface every version must satisfy.
// As long as a type has a GetLine method with this signature,
// it can be assigned to `active` below.
type Liner interface {
	GetLine(prompt string, buffer []byte) bool
}

// ─────────────────────────────────────────────────────────────────
// CHANGE ME — swap the struct literal to pick a different version.
// ─────────────────────────────────────────────────────────────────

//var active Liner = getlineV1{}

//var active Liner = getlineV2{}

//var active Liner = getlineV3{}

var active Liner = getlineV4{}

// ─────────────────────────────────────────────────────────────────

func main() {
	buffer := make([]byte, bufSize) // bufSize is defined in terminal.go

	// Pre-load a default value so Version Three has something to show.
	// Versions One and Two ignore any existing buffer contents.
	copy(buffer, "default text")

	fmt.Println("Get_Line demo — The Craft of Text Editing (Finseth)")
	fmt.Println()

	ok := active.GetLine("Enter some text", buffer)
	if !ok {
		fmt.Fprintln(os.Stderr, "GetLine returned false (cancelled or buffer too small)")
		os.Exit(1)
	}

	fmt.Printf("\nYou entered : %q\n", cstring(buffer))
	fmt.Printf("Length      : %d characters\n", len(cstring(buffer)))
}

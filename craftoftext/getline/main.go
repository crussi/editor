package main

import (
	"flag"
	"fmt"
	"os"
)

// Liner is the interface every version must satisfy.
type Liner interface {
	GetLine(prompt string, buffer []byte) bool
}

func main() {
	// Command-line flag to choose version
	version := flag.Int("v", 5, "GetLine version to use (1–5)")
	flag.Parse()

	var active Liner

	// Select version based on flag
	switch *version {
	case 1:
		active = getlineV1{}
	case 2:
		active = getlineV2{}
	case 3:
		active = getlineV3{}
	case 4:
		active = getlineV4{}
	case 5:
		active = getlineV5{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown version %d — using V5\n", *version)
		active = getlineV5{}
	}

	buffer := make([]byte, bufSize)

	// Pre-load a default value so Version Three+ has something to show.
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

package main

import (
	"flag"
	"fmt"
	"os"
)

type Liner interface {
	GetLine(prompt string, buffer []byte) bool
}

func main() {
	version := flag.Int("v", 5, "GetLine version to use (1–6)")
	flag.Parse()

	var active Liner

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
	case 6:
		active = getlineV6{
			Allowed: []string{"yes", "no", "maybe"},
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown version %d — using V5\n", *version)
		active = getlineV5{}
	}

	buffer := make([]byte, bufSize)

	// Only preload default text for versions 3–5
	switch *version {
	case 3, 4, 5:
		copy(buffer, "default text")
	}

	fmt.Println("Get_Line demo — The Craft of Text Editing (Finseth)")
	fmt.Println()

	// If this is getlineV6, show allowed responses
	if v6, ok := active.(getlineV6); ok {
		fmt.Printf("Allowed responses: %v\n", v6.Allowed)
	}

	ok := active.GetLine("Enter some text", buffer)
	if !ok {
		fmt.Fprintln(os.Stderr, "GetLine returned false (cancelled or buffer too small)")
		os.Exit(1)
	}

	fmt.Printf("\nYou entered : %q\n", cstring(buffer))
	fmt.Printf("Length      : %d characters\n", len(cstring(buffer)))
}

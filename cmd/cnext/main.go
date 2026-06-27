package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cnext <command>")
		fmt.Println("Commands:")
		fmt.Println("  init      Initialize a new cnext project")
		fmt.Println("  build     Build the project")
		fmt.Println("  test      Run tests")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		fmt.Println("cnext init - not yet implemented")
	case "build":
		fmt.Println("cnext build - not yet implemented")
	case "test":
		fmt.Println("cnext test - not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

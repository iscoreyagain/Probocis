package main

import (
	"fmt"
	"os"

	"github.com/iscoreyagain/Probocis/internals/commands"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
/*func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: prob <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
*/

func main() {
	// Ensure the command is invoked with "hash-object" as the first argument
	if len(os.Args) < 2 || os.Args[1] != "hash-object" {
		fmt.Fprintln(os.Stderr, "Usage: git hash-object [-w] [--stdin] [<file>]")
		os.Exit(1)
	}

	// Create an instance of HashObjCmd
	cmd := &commands.HashObjCmd{}

	// Pass arguments after "hash-object" to the Run method
	// os.Args[2:] skips the program name and "hash-object"
	err := cmd.Run(os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

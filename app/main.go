package main

import (
	"fmt"
	"os"

	"github.com/iscoreyagain/Probocis/internals/commands"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	if os.Args[1] == "shutdown" {
		fmt.Println("Shutting down gracefully...")
		return
	}

	cmdName := os.Args[1]
	cmd, err := commands.NewCommand(cmdName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := cmd.Run(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Probocis is a internal tool for managing version control like well-known Git, written in Go")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    probocis <command> [arguments]")
	fmt.Println()
	fmt.Println("Low-level Commands (Plumbing):")
	fmt.Println("    hash-object  compute object ID and optionally create a blob")
	fmt.Println("    cat-file  inspect repository objects")
	fmt.Println("    update-index  manipulate the index")
	fmt.Println("    write-tree  create a tree object from the the index")
	fmt.Println("    commit-tree  create a commit object")
	fmt.Println("    read-tree  populate the index from a tree object")
	fmt.Println("    update-ref  update object references")
	fmt.Println("    symbolic-ref  manage symbolic references")
	fmt.Println("High-level Commands (Porcelain):")
	fmt.Println("    init  initialize the repository")
	fmt.Println("    add  add files to the staging area")
	fmt.Println("    commit  create a new commit")
	fmt.Println("    status  display repository status")
	fmt.Println("    log  show commit history")
	fmt.Println("    diff  show file differences")
}

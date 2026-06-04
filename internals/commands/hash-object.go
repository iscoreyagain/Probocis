package commands

import (
	"fmt"

	"github.com/iscoreyagain/Probocis/internals/objects"
	"github.com/iscoreyagain/Probocis/internals/utils"
)

type HashObjCmd struct{}

func (h *HashObjCmd) Name() string {
	return "hash-object"
}

// git hash-object myfile.txt
// git hash-object -w myfile.txt
// echo "hello world" | git hash-object --stdin
// ["git", "hash-object", "-w", "myfile.txt", --stdin]
func (h *HashObjCmd) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: probocis hash-object [-w] [--stdin] <given_file>")
	}

	var isWrite bool
	var filename string
	var readFromStdin bool

	for _, arg := range args {
		switch arg {
		case "-w":
			isWrite = true
		case "--stdin":
			readFromStdin = true
		default:
			// Assuming the last non-flag argument is the filename
			filename = arg
		}
	}

	// Validation
	if filename != "" && readFromStdin {
		return fmt.Errorf("cannot use --stdin with a filename")
	} else if filename == "" && !readFromStdin {
		return fmt.Errorf("must specify a file or use --stdin")
	}

	var repoRoot string
	if isWrite {
		root, err := utils.FindRepoRoot()
		if err != nil {
			return fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
		}
		repoRoot = root
	}

	obj, err := objects.ReadObject(repoRoot, filename, readFromStdin)
	if err != nil {
		return fmt.Errorf("Failed to read object due to: %w", err)
	}

	objHash, objData, err := objects.HashObject(obj)
	if err != nil {
		return err
	}

	if isWrite {
		if err := objects.WriteObject(repoRoot, objHash, objData); err != nil {
			return err
		}
	}

	fmt.Print(objHash)
	return nil
}

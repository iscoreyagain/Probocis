package commands

import (
	"fmt"

	"github.com/iscoreyagain/Probocis/internals/objects"
)

type CatFileCmd struct{}

// Common structure: git cat-file <type> {[-p] [-t] [-s] [-e]} <object_hash_id>
func (c *CatFileCmd) Name() string {
	return "cat-file"
}

// git cat-file -t 213ea341ab567c8d...
func (c *CatFileCmd) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: git cat-file [-p|-t|-s] <sha>")
	}

	flags := args[0]
	hash_val := args[1]

	path, err := objects.GetPath(hash_val)

	object, err := objects.ReadGitObject(path)

	if err != nil {
		return fmt.Errorf("failed to exec the command due to: %w", err)
	}

	switch flags {
	case "-p":
		fmt.Print(string(object.Data()))
	case "-s":
		fmt.Println(len(object.Data()))
	case "-t":
		fmt.Println(object.Type())
	default:
		return fmt.Errorf("unknown flag: %s", flags)
	}

	return nil
}

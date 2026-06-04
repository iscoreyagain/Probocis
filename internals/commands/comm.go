package commands

import (
	"errors"
)

type Command interface {
	Name() string
	Run(args []string) error
}

func NewCommand(name string) (Command, error) {
	switch name {
	case "init":
		return &InitCmd{}, nil
	case "hash-object":
		return &HashObjCmd{}, nil
	case "cat-file":
		return &CatFileCmd{}, nil
	case "update-index":
		return &UpdateIndexCmd{}, nil
	case "ls-files":
		return &LsFilesCmd{}, nil
	default:
		return nil, errors.New("unknown or not supported command: " + name)
	}
}

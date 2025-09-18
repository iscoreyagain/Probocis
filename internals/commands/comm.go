package commands

import (
	"errors"
	"os"
)

type Command interface {
	Name() string
	Run(args []string) error
}

func ParseCmd() []string {
	args := os.Args[1:]

	return args
}

func NewCommand(name string) (Command, error) {
	switch name {
	case "init":
		return &InitCmd{}, nil
	case "clone":
		return &CloneCmd{}, nil
	case "commit":
		return &CommitCmd{}, nil
	default:
		return nil, errors.New("unknown or not supported command: " + name)
	}
}

package commands

import (
	"fmt"

	"github.com/iscoreyagain/Probocis/internals/objects"
)

type LsFilesCmd struct{}

func (l *LsFilesCmd) Name() string {
	return "ls-files"
}

func (l *LsFilesCmd) Run(args []string) error {
	entries, err := objects.LoadIndexFromDisk()
	if err != nil {
		return err
	}

	stage := false

	for _, arg := range args {
		if arg == "--stage" {
			stage = true
			break
		}
	}
	for _, e := range entries {
		if stage {
			fmt.Printf("%06o %x %d\t%s\n", e.Mode, e.Hash, e.Stage(), e.Path)
		} else {
			fmt.Println(e.Path)
		}
	}
	return nil
}

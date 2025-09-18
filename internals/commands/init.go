package commands

import (
	"fmt"
	"os"
)

type InitCmd struct{}

func (i *InitCmd) Name() string {
	return "init"
}

func (i *InitCmd) Run(args []string) error {
	//First, create the directory called ".git"
	err := os.Mkdir(".git", 0755)

	if err != nil {
		return fmt.Errorf("Failed to create .git directory: %w", err)
	}

	subDir := []string{"objects", "refs"}

	for _, dir := range subDir {
		err := os.MkdirAll(".git/"+dir, 0755)

		if err != nil {
			return fmt.Errorf("Failed to create %s: %w", dir, err)
		}
	}

	err = os.WriteFile(".git/HEAD", []byte("ref: refs/heads/master\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to create HEAD: %w", err)
	}

	fmt.Println("Initialized empty Git repository in .git/")
	return nil
}

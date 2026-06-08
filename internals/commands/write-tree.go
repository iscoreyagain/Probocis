package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/iscoreyagain/Probocis/internals/objects"
	"github.com/iscoreyagain/Probocis/internals/utils"
)

type WriteTreeCmd struct{}

func (w *WriteTreeCmd) Name() string {
	return "write-tree"
}

func (w *WriteTreeCmd) Run(args []string) error {
	repoRoot, err := utils.FindRepoRoot()
	if err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	root, err := objects.ConstructTreeFromEntries()
	if err != nil {
		return fmt.Errorf("failed to construct tree: %w", err)
	}

	hash, err := traversal(repoRoot, root)
	if err != nil {
		return fmt.Errorf("failed to write tree: %w", err)
	}

	fmt.Printf("%x\n", hash)
	return nil
}

func traversal(repoRoot string, node *objects.TreeNode) ([20]byte, error) {
	if node.Children == nil {
		return node.Hash, nil
	}

	var children []*objects.TreeNode

	for _, child := range node.Children {
		childHash, err := traversal(repoRoot, child)
		if err != nil {
			return [20]byte{}, err
		}

		child.Hash = childHash
		children = append(children, child)
	}

	tree := objects.NewTree(children)
	hexHash, data, err := objects.HashObject(tree)
	if err != nil {
		return [20]byte{}, err
	}

	err = objects.WriteObject(repoRoot, hexHash, data)
	if err != nil {
		return [20]byte{}, err
	}

	var hash [20]byte
	decoded, _ := hex.DecodeString(hexHash)
	copy(hash[:], decoded)
	return hash, nil
}

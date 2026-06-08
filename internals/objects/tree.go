package objects

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iscoreyagain/Probocis/internals/constants"
)

type Tree struct {
	entries []*TreeNode
}
type TreeNode struct {
	Name     string
	Hash     [20]byte
	Mode     constants.FileMode
	Children map[string]*TreeNode
}

func (t *Tree) Type() string {
	return "tree"
}

func (t *Tree) Data() []byte {
	var buf bytes.Buffer

	for _, entry := range t.entries {
		// Git tree entry format: "{mode} {name}\0{20-byte raw hash}"
		_, err := fmt.Fprintf(&buf, "%o %s\x00", entry.Mode, entry.Name)
		if err != nil {
			return nil
		}
		buf.Write(entry.Hash[:])
	}

	return buf.Bytes()
}

func NewTree(entries []*TreeNode) *Tree {
	return &Tree{entries: entries}
}

// ConstructTreeFromEntries loads the current index and builds an in-memory
// tree reflecting the staged directory structure.
func ConstructTreeFromEntries() (*TreeNode, error) {
	entries, err := LoadIndexFromDisk()
	if err != nil {
		return nil, err
	}

	root := &TreeNode{
		Children: make(map[string]*TreeNode),
	}

	for _, entry := range entries {
		parts := strings.Split(entry.Path, "/")
		curr := root

		for _, part := range parts[:len(parts)-1] {
			if _, ok := curr.Children[part]; !ok {
				curr.Children[part] = &TreeNode{
					Name:     part,
					Mode:     constants.ModeDirectory,
					Children: make(map[string]*TreeNode),
				}
			}
			curr = curr.Children[part]
		}
		name := parts[len(parts)-1]
		curr.Children[name] = &TreeNode{
			Name:     name,
			Mode:     entry.Mode,
			Hash:     entry.Hash,
			Children: nil,
		}
	}
	return root, nil
}

package objects

import "fmt"

type GitObject interface {
	Type() string
	Data() []byte
}

func NewGitObject(objType string, content []byte) (GitObject, error) {
	switch objType {
	case "blob":
		return NewBlob(content), nil

	case "tree":
		return NewTree(content), nil

	case "commit":
		return NewCommit(content), nil

	default:
		return nil, fmt.Errorf("unknown object type: %s", objType)
	}
}

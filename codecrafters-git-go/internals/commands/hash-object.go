package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type HashObjCmd struct{}

func (h *HashObjCmd) Name() string {
	return "hash-object"
}

// git hash-object myfile.txt
// git hash-object -w myfile.txt
// echo "hello world" | git hash-object --stdin

func (h *HashObjCmd) Run(args []string) error {
	obj, _ := ReadObject(args)
	var isWrite bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w":
			isWrite = true
		}
	}

	objHash, objData, err := HashObject(obj)
	if err != nil {
		return err
	}

	// Check whether we need to write contents to :
	// .git/objects/[2 characters of the hash]/[the rest characters (62 for sha256, 38 for sha1)]

	if isWrite {
		if err := WriteObject(objHash, objData); err != nil {
			return err
		}
	}
	return nil
}

func ReadObject(args []string) (GitObject, error) {
	var filename string
	var readFromStdin bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--stdin":
			readFromStdin = true
		default:
			filename = args[i]
		}
	}

	var data []byte
	var err error

	//Check whether it's reading from stdin or reading from an existing file
	if readFromStdin {
		var err error
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return GitObject{}, fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else if filename != "" {
		var err error
		data, err = os.ReadFile(filename)
		if err != nil {
			return GitObject{}, fmt.Errorf("failed to read from file: %w", err)
		}
	} else {
		return GitObject{}, err
	}

	return GitObject{Type: "blob", Content: data}, nil
}

func HashObject(git GitObject) (string, []byte, error) {
	if git.Type == "" {
		return "", nil, fmt.Errorf("invalid content")
	}
	//Append the header of blob object + the actual content of the file
	header := fmt.Sprintf("%s %d\x00", git.Type, len(git.Content))
	objData := append([]byte(header), git.Content...)

	//Hash the entire blob object
	objHash := sha256.Sum256([]byte(objData))

	return fmt.Sprintf("%x", objHash), objData, nil
}

func WriteObject(hash string, content []byte) error {
	first := hash[:2]
	rest := hash[2:]

	// Check if the DIRECTORY/FOLDER was existed before (if not creating a new one)
	path := fmt.Sprintf(".git/objects/%s", first)

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create object dir: %w", err)
	}

	// Compressed and write the compressed data to the buffer
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(content); err != nil {
		return fmt.Errorf("failed to compress the object content: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close the stream: %w", err)
	}

	// Write the compressed data from the previous buffer to the file (which is named from the remaining 62 hex chars)
	objPath := fmt.Sprintf("%s/%s", path, rest)

	//The file is read-only for everybody
	if err := os.WriteFile(objPath, b.Bytes(), 0444); err != nil {
		return fmt.Errorf("failed to write content to object file: %w", err)
	}

	return nil
}

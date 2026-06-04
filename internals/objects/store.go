package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iscoreyagain/Probocis/internals/utils"
)

// ** Pass the "absolute path" to the function and return:
// - type
// - size
// - actual "meaningful" content
func ReadGitObject(path string) (GitObject, error) {
	// Not decompressed yet
	whole_content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file due to: %w", err)
	}

	// Decompressing the file's content
	b := bytes.NewReader(whole_content)

	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read file due to: %w", err)
	}
	defer z.Close()

	decompressed, err := io.ReadAll(z)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Parse the entire file content
	parts := bytes.SplitN(decompressed, []byte{0}, 2)

	// check corruption
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid git object")
	}

	headerParts := strings.SplitN(string(parts[0]), " ", 2)

	// check corruption
	if len(headerParts) != 2 {
		return nil, fmt.Errorf("invalid object header")
	}

	return NewGitObject(headerParts[0], parts[1])
}

// Good tips/tricks: in the user's perspective, if they type something like this:
// git cat-file -t <some_hash_id> - do they want to read a file in their current working dir
// OR some absolute path? (for .i.e, /tmp/foo.txt)
func GetPath(hash string) (string, error) {
	if len(hash) != 40 {
		return "", fmt.Errorf("invalid hash length")
	}

	dir := hash[:2]
	file := hash[2:]

	repoRoot, err := utils.FindRepoRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(repoRoot, ".probocis", "objects", dir, file), nil
}

func ReadObject(repoRoot, filename string, readFromStdin bool) (GitObject, error) {
	var data []byte
	var err error

	//Check whether it's reading from stdin or reading from an existing file
	if readFromStdin {
		data, err = io.ReadAll(os.Stdin)
		fmt.Printf("%q\n", data)
		fmt.Printf("% x\n", data)
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		data, err = os.ReadFile(filename) //"blob 5\0hello"
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
	}

	fmt.Println(data) //debug

	return NewBlob(data), nil
}

func HashObject(git GitObject) (string, []byte, error) {
	if git == nil {
		return "", nil, fmt.Errorf("invalid git object")
	}
	//Append the header of blob object + the actual content of the file
	header := fmt.Sprintf("%s %d\x00", git.Type(), len(git.Data()))
	objData := append([]byte(header), git.Data()...)

	//Hash the entire blob object
	objHash := sha1.Sum(objData)

	return fmt.Sprintf("%x", objHash), objData, nil
}

func WriteObject(repoRoot, hash string, content []byte) error {
	first := hash[:2]
	rest := hash[2:]

	// Build the dir
	objDir := filepath.Join(repoRoot, ".probocis", "objects", first)

	if err := os.MkdirAll(objDir, 0755); err != nil { //0755 - permission mode
		return fmt.Errorf("failed to create object dir: %w", err)
	}

	// Compressed and write the compressed data to the buffer
	b, err := utils.Compress(content)
	if err != nil {
		return err
	}

	// Write the compressed data from the previous buffer to the file (which is named from the remaining 62 hex chars)
	objPath := filepath.Join(objDir, rest)

	//The file is read-only for everybody
	if err := os.WriteFile(objPath, b, 0444); err != nil {
		return fmt.Errorf("failed to write content to object file: %w", err)
	}

	return nil
}

package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/iscoreyagain/Probocis/internals/objects"
)

func FindRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")

		info, err := os.Stat(gitPath)
		if err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return "", os.ErrNotExist
}

func Compress(content []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)

	if _, err := w.Write(content); err != nil {
		return nil, fmt.Errorf("failed to compress the object content: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the stream: %w", err)
	}

	return buf.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)

	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read file due to: %w", err)
	}
	defer z.Close()

	decompressed, err := io.ReadAll(z)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	return decompressed, nil
}

// Resolve a relative path typed from users' CLI to absolute path
// example:
/* /home/
└── user/
    └── documents/
        └── my-project/
            ├── .git/
            ├── src/
            │   └── main.go
            └── README.md */
// User type and we get the full path: /home/user/documents/my-project/src/main.go
func ResolvePath(path string) (string, error) {
	return filepath.Abs(path)
}

func ComputeHash(content []byte) []byte {
	hash := sha1.Sum(content)

	return hash[:]
}

func ReadObject(repoRoot, filename string, readFromStdin bool) (objects.GitObject, error) {
	var data []byte
	var err error

	//Check whether it's reading from stdin or reading from an existing file
	if readFromStdin {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else if filename != "" {
		filePath := filepath.Join(repoRoot, filename) //"C:\Users\QUOC THAI\myfile.txt"
		data, err = os.ReadFile(filePath)             //"blob 5\0hello"
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no file specified!")
	}

	return objects.NewBlob(data), nil
}

func HashObject(git objects.GitObject) (string, []byte, error) {
	if git == nil {
		return "", nil, fmt.Errorf("invalid git object")
	}
	//Append the header of blob object + the actual content of the file
	header := fmt.Sprintf("%s %d\x00", git.Type(), len(git.Content()))
	objData := append([]byte(header), git.Content()...)

	//Hash the entire blob object
	objHash := sha1.Sum([]byte(objData))

	return fmt.Sprintf("%x", objHash), objData, nil
}

func WriteObject(repoRoot, hash string, content []byte) error {
	first := hash[:2]
	rest := hash[2:]

	// Build the dir
	objDir := filepath.Join(repoRoot, ".git", "objects", first)

	if err := os.MkdirAll(objDir, 0755); err != nil { //0755 - permission mode
		return fmt.Errorf("failed to create object dir: %w", err)
	}

	// Compressed and write the compressed data to the buffer
	b, err := Compress(content)
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

// Helper func to read ALL OF existing entries in .git/index
func ReadIndex(indexPath string) ([]objects.IndexEntry, error) {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil, err
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(content)

	var numEntries uint32
	binary.Read(reader, binary.BigEndian, &numEntries)

	var entries []objects.IndexEntry
	for i := 0; i < int(numEntries); i++ {
		var entry objects.IndexEntry
		// For path, we just read until meet a null terminator \x00
		if err := entry.Deserialize(reader); err != nil {
			return nil, fmt.Errorf("Failed to deserialize into human-readable because: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

/* func WriteIndex(indexPath string, content []byte) error {

} */

func GetNumEntries(indexPath string) (uint32, error) {
	file, err := os.Open(indexPath)
	if err != nil {
		return 0, fmt.Errorf("Failed to open index file!")
	}
	defer file.Close()

	var numEntries uint32

	err = binary.Read(file, binary.BigEndian, &numEntries)
	if err != nil {
		return 0, fmt.Errorf("Failed to read index file!")
	}

	return numEntries, nil
}

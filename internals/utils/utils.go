package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FindRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".probocis")

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
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, err = os.Lstat(abs); err != nil {
		return "", fmt.Errorf("path does not exist %s: %w", abs, err)
	}

	return abs, nil
}

func ComputeHash(content []byte) []byte {
	header := fmt.Sprintf("blob %d\x00", len(content))
	data := append([]byte(header), content...)
	hash := sha1.Sum(data)

	return hash[:]
}

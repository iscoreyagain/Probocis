package internals

import (
	"bytes"
	"fmt"
	"os"
)

func handleError(err error) {
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
}

// In Git, ALL objects are identifiable by a 40-character SHA-1 hash, also known as the "object hash".
// For .i.e: e88f7a929cd70b0274c4ea33b209c97fa845fdbc
// Git objects are stored in the .git/objects directory. The path to each object is derived from the hash value.
// .git/objects/e8/8f7a929cd70b0274c4ea33b209c97fa845fdbc

func readContentObject(hash string) (string, error) {
	if len(hash) != 40 {
		return "", fmt.Errorf("Invalid length of hash")
	}

}

func readTypeObject(hash string) (string, error) {

}

func readSizeObject(hash string) (int64, error) {

}

func readObject(hash string) bytes.Buffer {
	dir := fmt.Sprintf(".git/objects/%s", hash[:2])
	fileName := fmt.Sprintf("%s/%s", dir, hash[2:])

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error(s) while reading file: %v\n", err)
		os.Exit(1)
	}

	//Chưa xong
}

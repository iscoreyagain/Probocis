package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iscoreyagain/Probocis/internals/objects"
	"github.com/iscoreyagain/Probocis/internals/utils"
)

type IndexAction interface {
	Execute() error
}

type AddAction struct {
	RepoRoot, File string
}

func (a *AddAction) Execute() error {
	return add(a.RepoRoot, a.File)
}

type RemoveAction struct {
	RepoRoot, File string
}

func (a *RemoveAction) Execute() error {
	return remove(a.RepoRoot, a.File)
}

type UpdateIndexCmd struct{} //Class bên trong Java

func (u *UpdateIndexCmd) Name() string { //Method
	return "update-index"
}

// git update-index --add hello.txt
func (u *UpdateIndexCmd) Run(args []string) error {
	repoRoot, err := utils.FindRepoRoot() // ← thêm
	if err != nil {
		return err
	}

	for i := range len(args) {
		switch args[i] {
		case "--add":
			return (&AddAction{RepoRoot: repoRoot, File: args[i+1]}).Execute()
		case "--remove":
			return (&RemoveAction{RepoRoot: repoRoot, File: args[i+1]}).Execute()
		}
	}
	return fmt.Errorf("no valid flag provided")
}

// Besides ["git", "update-index"], it also currently supports essential flags:
// "--add"
// "--cacheinfo"
// "--remove"
// "--replace"
// "--refresh"

// git update-index --add src/hello.txt

// add stages a file into the index for the next commit.
// It resolves the given file path relative to repoRoot, computes its SHA-1 hash,
// and upserts the corresponding entry into .probocis/index.
// If the index file does not exist, it will be created automatically.
func add(repoRoot string, file string) error {
	// take the current working directory - CWD combine with the file
	filePath, err := utils.ResolvePath(file)
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(repoRoot, filePath)
	if err != nil {
		return err
	}

	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	hash := utils.ComputeHash(content)

	entry := objects.CreateNewEntry(relativePath, [20]byte(hash), fileInfo)

	entries, err := objects.LoadIndexFromDisk()
	if err != nil {
		return err
	}

	objects.Upsert(&entries, entry)

	return objects.SaveIndexToDisk(entries)
}

// git update-index --remove README.txt
// Purpose: Remove a file from the index (it will no longer be tracked)
func remove(repoRoot string, file string) error {
	return nil
}

// git update-index --cacheinfo 100644,5f6b8f...,file.txt
// Purpose: Add/update a specific object in the index by mode, SHA, and path directly
func cacheinfo() {

}

// git update-index --refresh
// Purpose: Update index entries with the current content of the working directory
func refresh() {

}

// git update-index --replace file.txt
// Purpose: Replace an existing file in the index
func replace() {

}

// Helper function
//func createEmptyIndex(indexPath string) error {
//	var numEntries uint32 = 0
//
//	index, err := os.Create(indexPath)
//	if err != nil {
//		return fmt.Errorf("Failed to create index file in .git repo: %w", err)
//	}
//
//	defer index.Close()
//
//	//Write number of total entries in the first 4 bytes
//	if err := binary.Write(index, binary.BigEndian, numEntries); err != nil {
//		return fmt.Errorf("Failed to initiate the number of entries due to: %w", err)
//	}
//
//	fmt.Println("Succesfully create empty index file at .git/index file with 0 entry")
//	return nil
//}

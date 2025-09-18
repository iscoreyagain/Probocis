package commands

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/git-starter-go/internals/objects"
	"github.com/codecrafters-io/git-starter-go/internals/utils"
)

type UpdateIndCmd struct{}

func (u *UpdateIndCmd) Name() string {
	return "update-index"
}

// git update-index --add hello.txt
func (u *UpdateIndCmd) Run(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("driver or directory not found!: %w", err)
	}

	var action string
	var filename string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--add":
			AddEntry
		}
	}
	//Get current working directory (~relative path)

	obj, err := ReadObject(args)
	if err != nil {
		return err
	}

	objHash
	if err := cmd.Run(args); err != nil {
		return err
	}

	return nil
}

// Besides ["git", "update-index"], it also currently supports essential flags:
// "--add"
// "--cacheinfo"
// "--remove"
// "--replace"
// "--refresh"

type UpdateOpts struct {
	Add       bool
	Remove    bool
	Replace   bool
	Refresh   bool
	CacheInfo *CacheInfoEntry
}

type CacheInfoEntry struct {
	Mode string
	Hash string
	Path string
}

// type ChmodEntry struct {Mode string, Path string}

// The func() will parse the users' command into Go structured-format to easily manipulate
func parseOptions(args []string) (*UpdateOpts, error) {
	opts := &UpdateOpts{}
	pos := 1

	for pos < len(args) {
		switch args[pos] {
		case "--add":
			pos++
			if pos >= len(args) {
				return nil, fmt.Errorf("missing path after --add")
			}
			opts.Add = true
		case "--remove":
			pos++
			if pos >= len(args) {
				return nil, fmt.Errorf("missing path after --remove")
			}
			opts.Remove = true
		case "--cacheinfo":
			pos++
			if pos >= len(args) {
				return nil, fmt.Errorf("missing cacheinfo value")
			}
			parts := strings.Split(args[pos], ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid cacheinfo format")
			}
			opts.CacheInfo = &CacheInfoEntry{
				Mode: parts[0], Hash: parts[1], Path: parts[2],
			}
		case "--replace":
			pos++
			opts.Replace = true
		case "--refresh":
			pos++
			opts.Refresh = true
		default:
			return nil, fmt.Errorf("unknown flag: %s", args[pos])
		}
		pos++
	}
	return opts, nil
}

// git update-index --add hello.txt
// Purpose: Add an existing file (entry) into .git/index for the next commit. In case if the index file not existed, create a new one then add
func add(repoRoot string, file string) error {
	indexPath := filepath.Join(repoRoot, ".git", "index")

	// check whether the current wd has .git/index yet
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if err := createEmptyIndex(indexPath); err != nil {
			return err
		}
	}

	//Read the existing file
	filePath, err := utils.ResolvePath(file)
	if err != nil {
		fmt.Errorf(filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	//Compute SHA-1
	hash := utils.ComputeHash(content)
	hashStr := fmt.Sprintf("%x", hash) //1

	//Create blob object in .git/objects
	if err := utils.WriteObject(repoRoot, hashStr, hash); err != nil {
		return err
	}

	//Adding it to .git/index file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to check file mode due to: %w", err)
	}

	fileMode := fileInfo.Mode()                           //2
	fileSize := fileInfo.Size()                           //3
	relativePath, err := filepath.Rel(repoRoot, filePath) //4
	if err != nil {
		return err
	}

	entry := &objects.IndexEntry{
		Mode: uint32(fileMode),
		Size: uint32(fileSize),
		Hash: [20]byte(hash),
		Path: relativePath,
	}

	n, err := getNumEntries(indexPath)
	if err != nil {
		return err
	}

	data, err := entry.Serialize()
	if err != nil {
		return err
	}

	if err := utils.WriteIndex(indexPath, data); err != nil {
		return fmt.Errorf("Failed to write/modify expected entry due to: %w", err)
	}

	return nil
}

// git update-index --remove README.txt
// Purpose: Remove a file from the index (it will no longer be tracked)
func remove(repoRoot string, file string) error {

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
func createEmptyIndex(indexPath string) error {
	var numEntries uint32 = 0

	index, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Failed to create index file in .git repo: %w", err)
	}

	defer index.Close()

	//Write number of total entries in the first 4 bytes
	if err := binary.Write(index, binary.BigEndian, numEntries); err != nil {
		return fmt.Errorf("Failed to initiate the number of entries due to: %w", err)
	}

	fmt.Println("Succesfully create empty index file at .git/index file with 0 entry")
	return nil
}

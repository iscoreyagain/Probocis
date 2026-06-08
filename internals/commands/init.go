package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iscoreyagain/Probocis/internals/constants"
	"github.com/iscoreyagain/Probocis/internals/objects"
	"github.com/iscoreyagain/Probocis/internals/utils"
)

type InitCmd struct{}

func (i *InitCmd) Name() string {
	return "init"
}

// example: [probocis init] --hash=sha256 --compress=zlib --verbose=true myrepo

// init creates or reinitialize the Probocis repo
//
// By default, the repository is created in the current working directory.
// If a target directory is provided, the repository is initialized there.
func (i *InitCmd) Run(args []string) error {
	conf := &objects.Config{
		Hash:          "sha1",
		Compress:      "zlib",
		Version:       0,
		DefaultBranch: "main",
	}

	var isVerbose = false

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--hash="):
			conf.Hash = strings.TrimPrefix(arg, "--hash=")
		case strings.HasPrefix(arg, "--compress="):
			conf.Compress = strings.TrimPrefix(arg, "--compress=")
		case arg == "--verbose":
			isVerbose = true
		}
	}

	logger := &utils.Logger{IsVerbose: isVerbose}

	repoDir, err := createRepo(parseTargetDir(args), logger)
	if err != nil {
		return err
	}

	if err = writeHead(repoDir, logger, conf); err != nil {
		return err
	}

	if err = writeConfig(repoDir, logger, conf); err != nil {
		return err
	}

	fmt.Printf(
		"Initialized empty Probocis repository in %s\n",
		repoDir,
	)

	return nil
}

// parseTargetDir return the directory specified by the user
// If no directory is provided, an empty string is returned.
func parseTargetDir(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			return arg
		}
	}

	return ""
}

// createRepo creates the Probocis repository structure in the target
// directory, including metadata files and required subdirectories.
func createRepo(targetDir string, logger *utils.Logger) (string, error) {
	defaultDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	defaultDir = filepath.Join(defaultDir, targetDir)

	repoDir := filepath.Join(defaultDir, constants.RepoDir)
	//First, create the directory called ".git"
	if _, err := os.Stat(repoDir); err == nil {
		return "", fmt.Errorf("repository already exists")
	}

	err = os.MkdirAll(repoDir, 0755)

	if err != nil {
		return "", fmt.Errorf("Failed to create .probocis directory: %w", err)
	}
	// log
	logger.Log("creating %s ", ".probocis/")

	dirs := []string{
		filepath.Join(repoDir, "objects"),
		filepath.Join(repoDir, "refs"),
		filepath.Join(repoDir, "refs", "heads")}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)

		if err != nil {
			return "", fmt.Errorf("Failed to create %s: %w", dir, err)
		}
		logger.Log("creating %s ", dir)
	}

	return repoDir, nil
}

// writeHead create the HEAD file and points it to the default branch
func writeHead(dir string, logger *utils.Logger, conf *objects.Config) error {
	headPath := filepath.Join(dir, constants.HeadFile)
	headContent := fmt.Sprintf("ref: refs/heads/%s\n", conf.DefaultBranch)
	err := os.WriteFile(headPath, []byte(headContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create HEAD: %w", err)
	}
	logger.Log("writing %s -> \n%s", headPath, headContent)

	return nil
}

// writeConfig writes to the repository metadata file, such as:
// - the hash algorithm
// - the compression method
// - the format version
func writeConfig(dir string, logger *utils.Logger, conf *objects.Config) error {
	configPath := filepath.Join(dir, "config")
	configContent := fmt.Sprintf(
		"[core]\n\t"+
			"repositoryformatversion = %d\n\t"+
			"bare = false\n\t"+
			"filemode = false\n\n"+
			"[probocis]\n\t"+
			"hash = %s\n\t"+
			"compress = %s\n\t"+
			"default_branch = %s\n",
		conf.Version,
		conf.Hash,
		conf.Compress,
		conf.DefaultBranch)
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create ./config: %w", err)
	}
	logger.Log("writing %s -> \n%s", configPath, configContent)

	return nil
}

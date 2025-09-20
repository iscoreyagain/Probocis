package commands

type CatFileCmd struct{}

// Common structure: git cat-file <type> {[-p] [-t] [-s] [-e]} <object_hash_id>
func (c *CatFileCmd) Name() string {
	return "cat-file"
}

// git cat-file -t 213ea341ab567c8d...
/* func (c *CatFileCmd) Run(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("usage: git cat-file [-p|-t|-s] <sha>")
	}

	flags := args[2]
	hash_id := args[3]

}

func normalizePath(input string) (string, error) {
	abs, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}

}

// Good tips/tricks: in the user's perspective, if they type something like this:
// git cat-file -t <some_hash_id> - do they want to read a file in their current working dir
// OR some absolute path? (for .i.e, /tmp/foo.txt)
func objectPath(hash string) string {
	dir := hash[:2]
	file := hash[2:]
	full_path := filepath.Join(".git", "objects", dir, file)

	return full_path
}

func readContent(path string) (string, error) {
	// Not decompressed yet
	whole_content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file due to: %w", err)
	}

	// Decompressing the file's content
	b := bytes.NewReader(whole_content)

	z, err := zlib.NewReader(b)
	if err != nil {
		return "", fmt.Errorf("failed to read file due to: %w", err)
	}
	defer z.Close()

	decompressed, err := io.ReadAll(z)
	if err != nil {
		return "", fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Parse the content
	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	content := parts[1]

	return string(content), nil
}
*/

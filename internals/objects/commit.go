package objects

// Define a simple Commit struct as a stub
type Commit struct {
	// Add fields as needed (e.g., Tree, Parent, Author, Committer, Message)
	data []byte
}

func (c *Commit) Type() string {
	return "commit"
}

func (c *Commit) Data() []byte {
	return c.data
}

// NewCommit creates a new Commit instance
func NewCommit(data []byte) *Commit {
	return &Commit{
		data: data,
	}
}

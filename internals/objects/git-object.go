package objects

type GitObject interface {
	Type() string
	Content() []byte
}

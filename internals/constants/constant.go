package constants

type FileMode uint32

const (
	ModeRegular    FileMode = 0o100644
	ModeExecutable FileMode = 0o100755
	ModeSymlink    FileMode = 0o120000
	ModeGitlink    FileMode = 0o160000
	ModeDirectory  FileMode = 0o040000
)

type ObjectID [20]byte

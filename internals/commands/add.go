package commands

type AddCmd struct{}

func (a *AddCmd) Name() string {
	return "add"
}

// git add hello.txt
// 1. Creates a blob object (hello.txt) in .git/objects
// 2. Create or updating BINARY FILE - .git/index
// 3. Add [hello.txt] entry in index file

func (a *AddCmd) Run(args []string) error {

}

func CreateIndexFile(objData)

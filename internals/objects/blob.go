package objects

type Blob struct {
	Data []byte
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Content() []byte {
	return b.Data
}

func NewBlob(data []byte) *Blob {
	return &Blob{Data: data}
}

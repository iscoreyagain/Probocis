package objects

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IndexEntry struct {
	Mode uint32
	Size uint32
	Hash [20]byte
	Path string
}

func (i *IndexEntry) Serialize() ([]byte, error) {
	var buff bytes.Buffer

	if err := binary.Write(&buff, binary.BigEndian, i.Mode); err != nil {
		return nil, err
	}
	if err := binary.Write(&buff, binary.BigEndian, i.Size); err != nil {
		return nil, err
	}
	if err := binary.Write(&buff, binary.BigEndian, i.Hash); err != nil {
		return nil, err
	}

	if _, err := buff.WriteString(i.Path); err != nil {
		return nil, err
	}

	padding := (8 - (buff.Len() % 8)) % 8
	if padding > 0 {
		if _, err := buff.Write(make([]byte, padding)); err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}

func (i *IndexEntry) Deserialize(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &i.Mode); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &i.Size); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &i.Hash); err != nil {
		return err
	}

	// Read the variable-length pathname
	var pathBuffer bytes.Buffer
	if br, ok := r.(io.ByteReader); ok {
		for {
			b, err := br.ReadByte()
			if err != nil {
				return err
			}
			if b == 0 {
				break
			}
			pathBuffer.WriteByte(b)
		}
	} else {
		// Fallback
		for {
			var b [1]byte
			if _, err := r.Read(b[:]); err != nil {
				return err
			}
			if b[0] == 0 {
				break
			}
			pathBuffer.WriteByte(b[0])
		}
	}

	i.Path = pathBuffer.String()

	totalSize := 28 + len(i.Path) + 1 // +1 for the null terminator
	padding := (8 - (totalSize % 8)) % 8

	if padding > 0 {
		if _, err := io.ReadFull(r, make([]byte, padding)); err != nil {
			return err
		}
	}

	return nil
}

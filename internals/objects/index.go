package objects

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/iscoreyagain/Probocis/internals/constants"
	"github.com/iscoreyagain/Probocis/internals/utils"
	_ "golang.org/x/sys/unix"
)

/*
	The full flow:

If the index file is not exist:

	CreateIndexFile() -> WriteHeader() -> Add() -> Save()

Else

	LoadIndexEntryFromDisk() -> cập nhật cái gì đó -> Save()
*/
type IndexEntry struct {
	CtimeSec  uint32 //4 bytes
	CtimeNSec uint32 //4 bytes
	MtimeSec  uint32 //4 bytes
	MtimeNSec uint32 //4 bytes
	Dev       uint32 //4 bytes
	Ino       uint32 //4 bytes
	UID       uint32 //4 bytes
	GID       uint32 //4 bytes
	Mode      constants.FileMode
	Size      uint32
	Hash      constants.ObjectID
	Flags     uint16
	Path      string
}

func (i *IndexEntry) Stage() uint16 {
	return (i.Flags >> 12) & 0x3
}

func ReadHeader(r io.Reader) (uint32, error) {
	var magic [4]byte
	var version, count uint32

	binary.Read(r, binary.BigEndian, &magic)
	if string(magic[:]) != "DIRC" {
		return 0, fmt.Errorf("invalid index file: bad magic %q", magic)
	}

	binary.Read(r, binary.BigEndian, &version)
	if version != 2 {
		return 0, fmt.Errorf("unsupported or invalid version: %d", version)
	}

	binary.Read(r, binary.BigEndian, &count)

	return count, nil
}

func writeHeader(w io.Writer, count uint32) error {
	err := binary.Write(w, binary.BigEndian, []byte("DIRC"))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, uint32(2))
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, count)
	if err != nil {
		return err
	}

	return nil
}

// ** The function doesn't need any kinds of arguments - it will always
// load the .probocis/index file from the current repo
func LoadIndexFromDisk() ([]*IndexEntry, error) {
	root, err := utils.FindRepoRoot()

	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(root, ".probocis", "index")

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*IndexEntry{}, nil // nho check
		}
		return nil, fmt.Errorf("failed to open index file from this path %s: %v", fullPath, err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	count, err := ReadHeader(reader)
	if err != nil {
		return nil, err
	}

	entries := make([]*IndexEntry, count)
	for i := range count {
		entry := &IndexEntry{}
		if err := entry.Deserialize(reader); err != nil {
			return nil, fmt.Errorf("failed to deserialize entry %d: %v", i, err)
		}
		entries[i] = entry
	}

	return entries, nil
}

func SaveIndexToDisk(entries []*IndexEntry) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	root, err := utils.FindRepoRoot()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(root, ".probocis", "index")

	var buf bytes.Buffer

	if err = writeHeader(&buf, uint32(len(entries))); err != nil {
		return err
	}

	for _, entry := range entries {
		data, err := entry.Serialize()
		if err != nil {
			return err
		}
		buf.Write(data)
	}

	checksum := sha1.Sum(buf.Bytes())
	buf.Write(checksum[:])

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

func CreateNewEntry(filepath string, hash [20]byte, info os.FileInfo) *IndexEntry {
	stat := info.Sys().(*syscall.Stat_t)
	return &IndexEntry{
		CtimeSec:  uint32(stat.Ctim.Sec),
		CtimeNSec: uint32(stat.Ctim.Nsec),
		MtimeSec:  uint32(stat.Mtim.Sec),
		MtimeNSec: uint32(stat.Mtim.Nsec),
		Dev:       uint32(stat.Dev),
		Ino:       uint32(stat.Ino),
		UID:       uint32(stat.Uid),
		GID:       uint32(stat.Gid),
		Mode:      normalizeMode(uint32(stat.Mode)),
		Size:      uint32(stat.Size),
		Hash:      hash,
		Flags:     uint16(len(filepath)) & 0x0FFF,
		Path:      filepath,
	}
}

func CreateIndexFile() error {
	root, err := utils.FindRepoRoot()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(root, ".probocis", "index")

	if _, err := os.Stat(fullPath); err == nil {
		return nil
	}

	var buf bytes.Buffer

	if err := writeHeader(&buf, 0); err != nil {
		return err
	}

	checksum := sha1.Sum(buf.Bytes())
	buf.Write(checksum[:])

	return os.WriteFile(fullPath, buf.Bytes(), 0644)
}

func (i *IndexEntry) Serialize() ([]byte, error) {
	var buff bytes.Buffer

	fixed := []any{
		i.CtimeSec, i.CtimeNSec,
		i.MtimeSec, i.MtimeNSec,
		i.Dev, i.Ino,
		i.Mode,
		i.UID, i.GID,
		i.Size,
		i.Hash,
		i.Flags,
	}
	for _, f := range fixed {
		if err := binary.Write(&buff, binary.BigEndian, f); err != nil {
			return nil, err
		}
	}

	buff.WriteString(i.Path)
	buff.WriteByte(0)

	length := 62 + len(i.Path) + 1

	if padding := 8 - (length % 8); padding < 8 {
		if _, err := buff.Write(make([]byte, padding)); err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}

// ** Deserialize the disk-based index entry from the given reader `r` -> the well-formatted in-memory IndexEntry struct
func (i *IndexEntry) Deserialize(r io.Reader) error {
	fixed := []any{
		&i.CtimeSec, &i.CtimeNSec,
		&i.MtimeSec, &i.MtimeNSec,
		&i.Dev, &i.Ino,
		&i.Mode,
		&i.UID, &i.GID,
		&i.Size,
		&i.Hash,
		&i.Flags,
	}
	for _, f := range fixed {
		if err := binary.Read(r, binary.BigEndian, f); err != nil {
			return err
		}
	}

	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	path, err := br.ReadString(0)
	if err != nil {
		return err
	}

	i.Path = strings.TrimRight(path, "\x00")

	entryLen := 62 + len(i.Path) + 1
	if pad := 8 - (entryLen % 8); pad < 8 {
		if _, err := io.ReadFull(br, make([]byte, pad)); err != nil {
			return err
		}
	}

	return nil
}

func normalizeMode(fMode uint32) constants.FileMode {
	switch {
	case fMode&syscall.S_IFLNK == syscall.S_IFLNK:
		return constants.ModeSymlink
	case fMode&syscall.S_IFDIR == syscall.S_IFDIR:
		return constants.ModeDirectory
	default:
		return constants.ModeRegular
	}
}

// ** Take the new entry and inserts if it doesn’t exist or updates an existing record if it does
func Upsert(entries *[]*IndexEntry, newEntry *IndexEntry) {
	for i, e := range *entries {
		if e.Path == newEntry.Path {
			(*entries)[i] = newEntry
			return
		}
	}
	*entries = append(*entries, newEntry)
}

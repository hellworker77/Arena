package util

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PakEntry struct {
	Name           string
	Offset         uint32
	Size           uint32
	CompressedSize uint32
}

func CreatePak(output string, files []string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	var entries []PakEntry
	dataBuf := &bytes.Buffer{}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		w.Write(content)
		w.Close()

		entries = append(entries, PakEntry{
			Name:           filepath.Base(f),
			Size:           uint32(len(content)),
			CompressedSize: uint32(compressed.Len()),
		})

		dataBuf.Write(compressed.Bytes())
	}
	headerSize := uint32(4 + 4)
	tableSize := uint32(0)

	for _, e := range entries {
		tableSize += 2
		tableSize += uint32(len(e.Name))
		tableSize += 4 + 4 + 4
	}

	offset := headerSize + tableSize

	currentOffset := offset
	for i := range entries {
		entries[i].Offset = currentOffset
		currentOffset += entries[i].CompressedSize
	}

	out.Write([]byte("PAK0"))

	binary.Write(out, binary.LittleEndian, uint32(len(entries)))

	for _, e := range entries {
		binary.Write(out, binary.LittleEndian, uint16(len(e.Name)))
		out.Write([]byte(e.Name))

		binary.Write(out, binary.LittleEndian, e.Offset)
		binary.Write(out, binary.LittleEndian, e.Size)
		binary.Write(out, binary.LittleEndian, e.CompressedSize)
	}

	out.Write(dataBuf.Bytes())

	return nil
}

func ReadPak(pakFile string) error {
	f, err := os.Open(pakFile)
	if err != nil {
		return err
	}
	defer f.Close()

	header := make([]byte, 4)
	if _, err := io.ReadFull(f, header); err != nil {
		return err
	}

	if string(header) != "PAK0" {
		return fmt.Errorf("invalid PAK file")
	}

	var count uint32
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return err
	}

	entries := make([]PakEntry, count)

	for i := uint32(0); i < count; i++ {

		var nameLen uint16
		if err := binary.Read(f, binary.LittleEndian, &nameLen); err != nil {
			return err
		}

		nameBytes := make([]byte, nameLen)
		if _, err := io.ReadFull(f, nameBytes); err != nil {
			return err
		}

		entries[i].Name = string(nameBytes)

		if err := binary.Read(f, binary.LittleEndian, &entries[i].Offset); err != nil {
			return err
		}
		if err := binary.Read(f, binary.LittleEndian, &entries[i].Size); err != nil {
			return err
		}
		if err := binary.Read(f, binary.LittleEndian, &entries[i].CompressedSize); err != nil {
			return err
		}
	}

	fmt.Println("PAK Contents:")
	for _, e := range entries {
		fmt.Printf(
			"Name: %-20s Offset: %-8d Size: %-8d Compressed: %-8d\n",
			e.Name, e.Offset, e.Size, e.CompressedSize,
		)
	}

	return nil
}

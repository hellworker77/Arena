package util

import (
	"os"
	"path/filepath"
)

type Property struct {
	Id   string `toml:"id"`
	Type string `toml:"type"`
}

type FileEntry struct {
	Name   string
	Offset uint64
	Size   uint64
}

func ArchiveAll() error {
	inputDir := "com"
	var files []string

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return CreatePak("com.pak", files)
}

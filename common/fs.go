package common

import (
	"fmt"
	"log"
	"os"
)

type WalkDirCallback func(path string, size int64)

type IFS interface {
	WalkDir(path string, callback WalkDirCallback) error
}

const DefaultMaxLevel = 4

type FS struct {
	maxLevel int
}

func (fs *FS) WalkDir(path string, callback WalkDirCallback) error {
	if fs.maxLevel == 0 {
		fs.maxLevel = DefaultMaxLevel
	}

	err := fs.walkDir(path, callback, 0)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FS) walkDir(path string, callback WalkDirCallback, level int) error {
	if level > fs.maxLevel {
		return nil
	}
	f, err := os.Open(path)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	if err != nil {
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("%s%s%s", path, os.PathSeparator, file.Name())
		if !file.IsDir() {
			callback(fmt.Sprintf("%s%s%s", path, "/", file.Name()), file.Size())
		} else {
			err := fs.walkDir(fmt.Sprintf("%s%s%s", path, os.PathSeparator, file.Name()), callback, level+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

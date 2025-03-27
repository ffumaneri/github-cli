package common

import (
	"fmt"
	"github.com/ffumaneri/github-cli/concurrency"
	"os"
	"time"
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
	wp := concurrency.NewWorkerPool(10)
	wp.Start()
	wp.AddTask(concurrency.Executor{func() error {
		return fs.walkDir(wp, path, callback, 0)
	}, func(err error) {
		fmt.Println(err)
	}})
	wp.WaitForTimeout(1 * time.Millisecond)
	return nil
}

func (fs *FS) walkDir(wp *concurrency.WorkerPool, path string, callback WalkDirCallback, level int) error {
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
	var children []string
	for _, file := range files {
		if !file.IsDir() {
			callback(fmt.Sprintf("%s%s%s", path, "/", file.Name()), file.Size())
		} else {
			children = append(children, fmt.Sprintf("%s%s%s", path, "/", file.Name()))
		}
	}
	for _, child := range children {
		wp.AddTask(concurrency.Executor{func() error {
			return fs.walkDir(wp, child, callback, level+1)
		}, func(err error) {
			fmt.Println(err)
		}})
	}

	return nil
}

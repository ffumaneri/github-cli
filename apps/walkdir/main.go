package main

import (
	"fmt"
	"github.com/ffumaneri/github-cli/common"
	"time"
)

// Main function para testing
func main() {
	ws := &common.FS{}
	t := time.Now()
	err := ws.WalkDir("/Users/facundofumaneri/personal", func(path string, size int64) {
		fmt.Printf("%s %d\n", path, size)
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("time: %s\n", time.Since(t))
}

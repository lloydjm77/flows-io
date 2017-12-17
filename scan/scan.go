package main

import (
	"fmt"
	"os"

	"github.com/lloydjm77/flows-io/fileutil"
)

func main() {
	home := os.Getenv("HOME")
	files := fileutil.FindSourceFiles(home + "/software-development/strategy-generator")
	for _, file := range files {
		fmt.Println(file)
	}
}

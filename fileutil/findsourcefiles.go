package fileutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fatih/set.v0"
)

const textType = "text/plain"

var excludedDirectories = set.NewNonTS(".git", ".settings", ".vscode", "target")

func FindSourceFiles(dir string) []string {
	var files []string

	err := filepath.Walk(dir, visit(&files))
	if err != nil {
		panic(err)
	}

	return files
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() && excludedDirectories.Has(info.Name()) {
			fmt.Printf("Found exclusion: %v - skipping.\n", info.Name())
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.Contains(DetectFileType(path), textType) {
			*files = append(*files, path)
		}
		return nil
	}
}

package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fatih/set.v0"
)

const textType = "text/"

var excludedDirectories = set.NewNonTS(".git", ".settings", ".vscode", "target")

func FindSourceFiles(dir string) []string {
	var files []string

	filepath.Walk(dir, visit(&files))

	return files
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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

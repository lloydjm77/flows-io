package fileutil

import (
	"bufio"
	"os"

	"github.com/lloydjm77/flows-io/analyzer"
)

var analyzers = []analyzer.StringAnalyzer{analyzer.HTTPStringAnalyzer{}}

// AnalyzeFile analyzes a file at the supplied location. A slice of
// analyzer.AnalysisResult structs is returned for matching elements.
func AnalyzeFile(path string) []analyzer.AnalysisResult {
	file, _ := os.Open(path)
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	analysisResults := []analyzer.AnalysisResult{}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		for _, analyzer := range analyzers {
			analysisResult := analyzer.Analyze(line)
			if analysisResult.Type == "" {
				continue
			}
			analysisResult.Path = path
			analysisResults = append(analysisResults, analysisResult)
		}
	}

	return analysisResults
}

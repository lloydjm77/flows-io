package fileutil

import (
	"bufio"
	"os"
	"sync"

	"github.com/lloydjm77/flows-io/analyzer"
)

var analyzers = []analyzer.StringAnalyzer{analyzer.HTTPStringAnalyzer{}}

func AnalyzeFile(path string, results chan analyzer.AnalysisResult, wg *sync.WaitGroup, inprogress chan int) {
	defer wg.Done()
	defer func(inprogressLocal chan int) {
		<-inprogressLocal
	}(inprogress)

	file, _ := os.Open(path)
	defer file.Close()

	inprogress <- 1

	fileScanner := bufio.NewScanner(file)

	// fmt.Printf("Analyzing file: %v\n", path)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		for _, analyzer := range analyzers {
			analysisResult := analyzer.Analyze(line)
			if analysisResult.Type == "" {
				continue
			}
			analysisResult.Path = path
			results <- analysisResult
		}
	}
}

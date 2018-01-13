package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/lloydjm77/flows-io/analyzer"
	"github.com/lloydjm77/flows-io/worker"

	"github.com/lloydjm77/flows-io/fileutil"
)

func run() {
	start := time.Now()

	args := os.Args[1:]
	files := fileutil.FindSourceFiles(args[0])

	results := []analyzer.AnalysisResult{}

	for _, file := range files {
		results = append(results, fileutil.AnalyzeFile(file)...)
	}

	for _, result := range results {
		fmt.Printf("AnalysisResult: %v\n", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nFound %v potential matches across %v files in %s.\n", len(results)+1, len(files)+1, elapsed)
}

func runConcurrent(maxConcurrent int) {
	start := time.Now()

	args := os.Args[1:]
	files := fileutil.FindSourceFiles(args[0])

	tasks := []worker.Task{}

	for _, file := range files {
		tasks = append(tasks, createTask(file))
	}

	// Use context to cancel goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChannel := worker.PerformTasks(ctx, tasks, maxConcurrent)

	// Print value from first goroutine and cancel others
	count := 1
	for result := range resultChannel {
		analysisResults := result.([]analyzer.AnalysisResult)
		for _, analysisResult := range analysisResults {
			fmt.Printf("AnalysisResult: %v\n", analysisResult)
			count++
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("\nFound %v potential matches across %v files in %s.\n", count, len(files)+1, elapsed)
}

func createTask(file string) worker.Task {
	return func() interface{} {
		return fileutil.AnalyzeFile(file)
	}
}

func main() {
	runConcurrent(runtime.NumCPU() / 2)
}

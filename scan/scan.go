package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lloydjm77/flows-io/analyzer"

	"github.com/lloydjm77/flows-io/fileutil"
)

func main() {
	start := time.Now()

	results := make(chan analyzer.AnalysisResult)
	spawnRoutines(results)

	elapsed := time.Since(start)
	fmt.Printf("\nTime elapsed: %s.\n", elapsed)
}

func spawnRoutines(results chan analyzer.AnalysisResult) {
	args := os.Args[1:]
	files := fileutil.FindSourceFiles(args[0])
	var wg sync.WaitGroup

	filesAnalyzed := len(files)
	wg.Add(filesAnalyzed)

	inprogress := make(chan int, 10)

	for _, file := range files {
		go fileutil.AnalyzeFile(file, results, &wg, inprogress)
	}

	go func() {
		fmt.Println("waiting")
		wg.Wait()
		fmt.Println("closing")
		close(results)
	}()

	count := 0
	for result := range results {
		fmt.Printf("Type: %v, Path: %v, Value: %v\n", result.Type, result.Path, result.Value)
		count++
	}
	// time.Sleep(10 * time.Second)
	// var a []analyzer.AnalysisResult

	// for i := 0; i < filesAnalyzed; i++ {
	// 	select {
	// 	case result := <-results:
	// 		fmt.Printf("\nFound %v potential matches.\n", len(a))
	// 		fmt.Printf("receiving %v\n", result)
	// 		a = append(a, result)
	// 	}
	// }

	// for _, result := range a {
	// 	fmt.Printf("Type: %v, Path: %v, Value: %v\n", result.Type, result.Path, result.Value)
	// }

	fmt.Printf("\nAnalyzed %v files.", filesAnalyzed)
	// fmt.Printf("\nFound %v potential matches.\n", len(a))
	fmt.Printf("\nFound %v potential matches.\n", count)
}

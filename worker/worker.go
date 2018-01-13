package worker

import (
	"context"
	"sync"
)

// Task is a function type for tasks to be performed.
// All incoming tasks have to conform to this function type.
type Task func() interface{}

// PerformTasks is a function which will be called by the client to perform
// multiple task concurrently.
// Input:
// tasks: the slice with functions (type TaskFunction)
// ctx:  the context to monitor to trigger the end of task processing and return
// maxConcurrent: the maxmimum number of concurrent goroutines
// Output: the channel with results
func PerformTasks(ctx context.Context, tasks []Task, maxConcurrent int) chan interface{} {
	buffer := make(chan struct{}, maxConcurrent)

	// Create a worker for each incoming task
	workers := make([]chan interface{}, 0, len(tasks))

	for _, task := range tasks {
		resultChannel := newWorker(ctx, task, buffer)
		workers = append(workers, resultChannel)
	}

	// Merge results from all workers
	out := merge(ctx, workers)
	return out
}

func newWorker(ctx context.Context, task Task, buffer chan struct{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		buffer <- struct{}{}
		defer close(out)

		select {
		case <-ctx.Done():
			// Received a signal to abandon further processing
			return
		case out <- runTask(task, buffer):
			// Got some result
		}
	}()

	return out
}

func runTask(task Task, buffer chan struct{}) interface{} {
	result := task()
	<-buffer
	return result
}

func merge(ctx context.Context, workers []chan interface{}) chan interface{} {
	// Merged channel with results
	out := make(chan interface{})

	// Synchronization over channels: do not close "out" before all tasks are completed
	var wg sync.WaitGroup

	// Define function which waits the result from worker channel
	// and sends this result to the merged channel.
	// Then it decreases the counter of running tasks via wg.Done().
	output := func(c <-chan interface{}) {
		defer wg.Done()
		for result := range c {
			select {
			case <-ctx.Done():
				// Received a signal to abandon further processing
				return
			case out <- result:
				// some message or nothing
			}
		}
	}

	wg.Add(len(workers))
	for _, workerChannel := range workers {
		go output(workerChannel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

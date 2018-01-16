package worker

import (
	"context"
	"sync"
)

// Task is a function type for tasks to be performed.
// All incoming tasks have to conform to this function type.
type Task func() interface{}

// IWaitGroup defines an interface for the sync.WaitGroup struct
// allowing us to mock values for testing.
type IWaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}

// TaskContext defines the context for task execution.
type TaskContext struct {
	Context       context.Context
	MaxConcurrent int
}

type workerContext struct {
	ctx    context.Context
	task   Task
	buffer chan struct{}
}

// PerformTasks is a function which will be called by the client to perform
// multiple task concurrently.
// Input:
// tasks: the slice with functions (type TaskFunction)
// ctx:  the context to monitor to trigger the end of task processing and return
// maxConcurrent: the maxmimum number of concurrent goroutines
// Output: the channel with results
func PerformTasks(taskContext *TaskContext, tasks []Task) chan interface{} {
	buffer := make(chan struct{}, taskContext.MaxConcurrent)

	// Create a worker for each incoming task
	workers := make([]chan interface{}, 0, len(tasks))

	for _, task := range tasks {
		resultChannel := newWorker(&workerContext{taskContext.Context, task, buffer})
		workers = append(workers, resultChannel)
	}

	// Merge results from all workers
	out := merge(taskContext.Context, workers, newWaitGroup())
	return out
}

func newWorker(workerContext *workerContext) chan interface{} {
	out := make(chan interface{})
	go runTask(workerContext, out)
	return out
}

func runTask(workerContext *workerContext, out chan interface{}) {
	workerContext.buffer <- struct{}{}
	defer close(out)

	select {
	case <-workerContext.ctx.Done():
		// Received a signal to abandon further processing
		return
	case out <- executeTask(workerContext.task, workerContext.buffer):
		// Got some result
	}
}

func executeTask(task Task, buffer chan struct{}) interface{} {
	result := task()
	<-buffer
	return result
}

func merge(ctx context.Context, workers []chan interface{}, wg IWaitGroup) chan interface{} {
	// Merged channel with results
	out := make(chan interface{})

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

func newWaitGroup() IWaitGroup {
	return new(sync.WaitGroup)
}

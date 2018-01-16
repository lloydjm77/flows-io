package worker

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformTasks(t *testing.T) {
	type args struct {
		taskContext TaskContext
		tasks       []Task
	}
	tests := []struct {
		name string
		args args
		want chan interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PerformTasks(&tt.args.taskContext, tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PerformTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newWorker(t *testing.T) {
	type args struct {
		workerContext workerContext
	}
	tests := []struct {
		name   string
		args   args
		cancel bool
		want   interface{}
	}{
		{name: "should_create_a_new_worker", args: args{workerContext{task: func() interface{} { return 1 }}}, cancel: false, want: 1},
		{name: "should_create_a_new_worker_that_gets_canceled", args: args{workerContext{task: func() interface{} { return 1 }}}, cancel: true},
	}
	for _, tt := range tests {
		// Set up the context.
		ctx, cancel := context.WithCancel(context.Background())
		tt.args.workerContext.ctx = ctx

		// Set up the buffer.
		buffer := make(chan struct{}, 1)
		tt.args.workerContext.buffer = buffer

		t.Run(tt.name, func(t *testing.T) {
			// Execute and get output channel.
			got := newWorker(&tt.args.workerContext)

			assert.Len(t, buffer, 0)
			if tt.cancel {
				cancel()
				assert.Len(t, got, 0)
			} else {
				assert.ObjectsAreEqual(tt.want, <-got)
			}
		})
	}
}

func Test_newWaitGroup(t *testing.T) {
	t.Run("should_return_a_new_wait_group", func(t *testing.T) {
		wg := newWaitGroup()
		assert.IsType(t, new(sync.WaitGroup), wg)
		assert.NotNil(t, wg)
	})
}

func Test_executeTask(t *testing.T) {
	type args struct {
		task   Task
		buffer chan struct{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "should_execute_task", args: args{task: func() interface{} { return 1 }}, want: 1},
	}
	for _, tt := range tests {
		buffer := make(chan struct{}, 1)
		buffer <- struct{}{}
		tt.args.buffer = buffer
		t.Run(tt.name, func(t *testing.T) {
			got := executeTask(tt.args.task, tt.args.buffer)
			assert.ObjectsAreEqual(tt.want, got)
			assert.Len(t, buffer, 0)
		})
		close(buffer)
	}
}

func Test_merge(t *testing.T) {
	type args struct {
		ctx     context.Context
		results [][]interface{}
	}
	tests := []struct {
		name   string
		args   args
		cancel bool
		want   []interface{}
	}{
		{
			name: "should_merge_results",
			args: args{
				results: [][]interface{}{
					[]interface{}{1, 2}, []interface{}{3, 4},
				}},
			cancel: false,
			want:   []interface{}{1, 2, 3, 4},
		},
		{
			name: "should_cancel_and_return",
			args: args{
				results: [][]interface{}{
					[]interface{}{1, 2}, []interface{}{3, 4},
				}},
			cancel: true,
			want:   []interface{}{},
		},
	}
	for _, tt := range tests {
		workers := make([]chan interface{}, 0, len(tt.args.results))
		for _, result := range tt.args.results {
			var wg sync.WaitGroup
			wg.Add(len(result))

			input := make(chan interface{})
			for _, subresult := range result {
				go func(val interface{}) {
					defer wg.Done()
					input <- val
				}(subresult)
			}

			workers = append(workers, input)

			go func() {
				wg.Wait()
				close(input)
			}()
		}

		// Set up the context.
		ctx, cancel := context.WithCancel(context.Background())
		tt.args.ctx = ctx

		t.Run(tt.name, func(t *testing.T) {
			count := 0
			results := merge(tt.args.ctx, workers, newFakeWaitGroup())
			if tt.cancel {
				cancel()
				assert.Empty(t, results)
			} else {
				for result := range results {
					assert.Contains(t, tt.want, result)
					count++
				}
			}
			assert.Equal(t, len(tt.want), count)
		})
	}
}

func newFakeWaitGroup() IWaitGroup {
	var fwg FakeWaitGroup
	fwg.wg = sync.WaitGroup{}
	return &fwg
}

type IWaitGroupCallTracker struct {
	addCount  int
	doneCount int
	waitCount int
}

type FakeWaitGroup struct {
	wg sync.WaitGroup
}

func (fwg *FakeWaitGroup) Add(delta int) {
	fmt.Println("fake add")
	fwg.wg.Add(delta)
}

func (fwg *FakeWaitGroup) Done() {
	fmt.Println("fake done")
	fwg.wg.Done()
}

func (fwg *FakeWaitGroup) Wait() {
	fmt.Println("fake wait")
	fwg.wg.Wait()
}

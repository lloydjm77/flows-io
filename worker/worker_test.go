package worker

import (
	"context"
	"reflect"
	"testing"
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

			if tt.cancel {
				cancel()
				if len(got) == 0 && len(buffer) == 0 {
					// Do nothing.
				} else {
					t.Error("newWorker() should have been canceled")
				}
			} else {
				// Receive value from channel to assert.
				val := <-got
				if reflect.DeepEqual(val, tt.want) && len(buffer) == 0 {
					// Do nothing.
				} else {
					t.Errorf("newWorker() = %v, want %v", val, tt.want)
				}
			}
		})
	}
}

func Test_runTask(t *testing.T) {
	type args struct {
		workerContext workerContext
		out           chan interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTask(&tt.args.workerContext, tt.args.out)
		})
	}
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
			if !reflect.DeepEqual(got, tt.want) || len(buffer) != 0 {
				t.Errorf("executeTask() = %v, want %v", got, tt.want)
			}
		})
		close(buffer)
	}
}

func Test_merge(t *testing.T) {
	type args struct {
		ctx     context.Context
		workers []chan interface{}
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
			if got := merge(tt.args.ctx, tt.args.workers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

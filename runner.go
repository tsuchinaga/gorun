package gorun

import (
	"context"
	"errors"
	"sync"
	"time"
)

func NewRunner() Runner {
	return &runner{tasks: []runnableTask{}}
}

var (
	NoTasksError = errors.New("runner has no tasks")
)

type Runner interface {
	AddTask(task Task)             // タスクの追加
	Run(ctx context.Context) error // 開始
}

type runner struct {
	tasks []runnableTask
	mx    sync.Mutex
}

func (r *runner) AddTask(task Task) {
	if task == nil {
		return
	}

	r.mx.Lock()
	defer r.mx.Unlock()
	r.tasks = append(r.tasks, &taskRunner{task: task})
}

func (r *runner) Run(ctx context.Context) error {
	if err := r.start(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (r *runner) start(ctx context.Context) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if len(r.tasks) == 0 {
		return NoTasksError
	}

	for _, t := range r.tasks {
		go t.run(ctx)
	}
	return nil
}

type runnableTask interface {
	run(ctx context.Context)
}

type taskRunner struct {
	task Task
}

func (r *taskRunner) run(ctx context.Context) {
	for {
		now := time.Now()
		d := r.task.NextTime(now)
		if d <= 0 {
			return
		}

		select {
		case <-time.After(d):
			go r.task.Run(ctx)
		case <-ctx.Done():
			return
		}
	}
}

package gorun

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func Test_NewRunner(t *testing.T) {
	t.Parallel()
	want := &runner{tasks: []runnableTask{}}
	got := NewRunner()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), want, got)
	}
}

func Test_runner_AddTask(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		arg  Task
		want []runnableTask
	}{
		{name: "引数がnilなら何もしない", arg: nil, want: []runnableTask{}},
		{name: "引数がTaskならtasksに追加される", arg: &testTask{}, want: []runnableTask{&taskRunner{task: &testTask{}}}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &runner{tasks: []runnableTask{}}
			runner.AddTask(test.arg)
			got := runner.tasks
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_runner_start(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		tasks []runnableTask
		want  error
	}{
		{name: "tasksが空ならerror", tasks: []runnableTask{}, want: NoTasksError},
		{name: "tasksが空じゃなければ実行してnilを返す",
			tasks: []runnableTask{&testRunnableTask{}},
			want:  nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runner := &runner{tasks: test.tasks}
			ctx, cancel := context.WithCancel(context.Background())
			got := runner.start(ctx)
			cancel()
			if !errors.Is(got, test.want) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_runner_Run(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		tasks []runnableTask
		want  error
	}{
		{name: "tasksが空ならerror", tasks: []runnableTask{}, want: NoTasksError},
		{name: "tasksが空じゃなければ実行してnilを返す",
			tasks: []runnableTask{&testRunnableTask{}},
			want:  nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			runner := &runner{tasks: test.tasks}
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			got := runner.Run(ctx)
			cancel()
			if !errors.Is(got, test.want) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

func Test_taskRunner_run(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		task *testTask
		want int
	}{
		{name: "次回実行時刻が過去なら終了する", task: &testTask{addTime: -1}, want: 0},
		{name: "次回実行時刻が現在なら終了する", task: &testTask{addTime: 0}, want: 0},
		{name: "contextがDoneされるまで繰り返される", task: &testTask{addTime: 100 * time.Millisecond}, want: 9},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			task := taskRunner{task: test.task}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			task.run(ctx)
			cancel()
			got := test.task.count
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

type testTask struct {
	Task
	addTime time.Duration
	count   int
}

func (t *testTask) NextTime(now time.Time) time.Time {
	return now.Add(t.addTime)
}
func (t *testTask) Run(context.Context) {
	t.count++
}

type testRunnableTask struct {
	runnableTask
}

func (t *testRunnableTask) run(ctx context.Context) {
	<-ctx.Done()
}

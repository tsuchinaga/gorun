package gorun

import (
	"reflect"
	"testing"
)

func Test_Name(t *testing.T) {
	t.Parallel()
	want := &runner{tasks: []taskRunner{}}
	got := NewRunner()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), want, got)
	}
}

type testTask struct {
	Task
}

func Test_runner_AddTask(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		arg  Task
		want []taskRunner
	}{
		{name: "引数がnilなら何もしない", arg: nil, want: []taskRunner{}},
		{name: "引数がTaskならtasksに追加される", arg: testTask{}, want: []taskRunner{{task: testTask{}}}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &runner{tasks: []taskRunner{}}
			runner.AddTask(test.arg)
			got := runner.tasks
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want, got)
			}
		})
	}
}

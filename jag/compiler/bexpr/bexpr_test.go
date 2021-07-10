package bexpr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func newTaskWithNeeds(name, needs string) *v1alpha1.PipelineTask {
	return &v1alpha1.PipelineTask{
		Name:       name,
		Task:       &v1alpha1.TaskSpec{Needs: needs},
		Dependants: []*v1alpha1.PipelineTask{},
	}
}

func TestEvalNeedsExpr(t *testing.T) {
	statuses := map[string]*v1alpha1.TaskStatus{
		"task1": {State: v1alpha1.TaskState_TASK_STATE_SUCCESS},
		"task2": {State: v1alpha1.TaskState_TASK_STATE_FAILURE},
		"task3": {State: v1alpha1.TaskState_TASK_STATE_PENDING},
		"task4": {State: v1alpha1.TaskState_TASK_STATE_ERROR},
		"task5": {State: v1alpha1.TaskState_TASK_STATE_SUCCESS},
	}

	tests := map[string]struct {
		task   *v1alpha1.PipelineTask
		want   bool
		errMsg string
	}{
		"Empty should be ready": {
			task: newTaskWithNeeds("", ""),
			want: true,
		},
		"Circular dependency error": {
			task:   newTaskWithNeeds("task1", "task1"),
			errMsg: "circular dependency in task: task1",
		},
		"Unary needs success ready": {
			task: newTaskWithNeeds("", "task1"),
			want: true,
		},
		"Unary needs success not ready": {
			task: newTaskWithNeeds("", "task2"),
			want: false,
		},
		"Unary needs failure ready": {
			task: newTaskWithNeeds("", "!task2"),
			want: true,
		},
		"Unary needs failure not ready": {
			task: newTaskWithNeeds("", "!task1"),
			want: false,
		},
		"Unary needs completed ready (success)": {
			task: newTaskWithNeeds("", "^task1"),
			want: true,
		},
		"Unary needs completed ready (failure)": {
			task: newTaskWithNeeds("", "^task2"),
			want: true,
		},
		"Unary needs completed not ready (pending)": {
			task: newTaskWithNeeds("", "^task3"),
			want: false,
		},
		"Binary And ready (success && success)": {
			task: newTaskWithNeeds("", "task1 && task5"),
			want: true,
		},
		"Binary And not ready (success && failure)": {
			task: newTaskWithNeeds("", "task1 && task2"),
			want: false,
		},
		"Binary Or ready (success && failure)": {
			task: newTaskWithNeeds("", "task1 || task2"),
			want: true,
		},
		"Binary Or ready (success && success)": {
			task: newTaskWithNeeds("", "task1 || task5"),
			want: true,
		},
		"Complex 1": {
			task: newTaskWithNeeds("", "task1 && !(task2 || task3)"),
			want: true,
		},
		"Complex 2": {
			task: newTaskWithNeeds("", "task2 || task3 || task4"),
			want: false,
		},
		"Complex 3": {
			task: newTaskWithNeeds("", "task2 || ^task3 || task4"),
			want: false,
		},
		"Complex 4": {
			task: newTaskWithNeeds("", "task1 && task5 && task1"),
			want: true,
		},
		"Complex 5": {
			task: newTaskWithNeeds("", "(task1 && task2) || (task1 && ^task2)"),
			want: true,
		},
		"Complex 6": {
			task: newTaskWithNeeds("", "!((task1 && task2) || (task1 && ^task2))"),
			want: false,
		},
		"Err 1": {
			task:   newTaskWithNeeds("", "t1"),
			errMsg: "bad identifier: t1",
		},
		"Err 2": {
			task:   newTaskWithNeeds("", "task1 && "),
			errMsg: "no match found",
		},
	}
	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			got, err := EvalNeedsExpr(tt.task, statuses)
			if tt.errMsg != "" {
				if err == nil {
					t.Errorf("expected an error to be thrown with: %s", tt.errMsg)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected an error: %s to contain: %s", err, tt.errMsg)
					return
				}
			}
			if got != tt.want {
				t.Errorf("EvalNeedsExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _result bool

func BenchmarkEvalNeedsExpr(b *testing.B) {
	statuses := map[string]*v1alpha1.TaskStatus{
		"task1": {State: v1alpha1.TaskState_TASK_STATE_SUCCESS},
		"task2": {State: v1alpha1.TaskState_TASK_STATE_FAILURE},
		"task3": {State: v1alpha1.TaskState_TASK_STATE_PENDING},
		"task4": {State: v1alpha1.TaskState_TASK_STATE_ERROR},
		"task5": {State: v1alpha1.TaskState_TASK_STATE_SUCCESS},
	}

	benchExpr := []string{
		"task1 && task5 && task1",
		"!((task1 && task2) || (task1 && ^task2))",
		"!((task1 && task2) || (task1 && ^task2)) && ((task1 && task2) || (task1 && ^task2))",
		"!((task1 && task2) || (task1 && ^task2)) && ((task1 && task2) || (task1 && ^task2)) || !((task1 && task2) || (task1 && ^task2))",
	}

	for i, expr := range benchExpr {
		b.Run(fmt.Sprintf("expr-%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_r, err := EvalNeedsExpr(newTaskWithNeeds("", expr), statuses)
				if err != nil {
					b.Errorf("unexpected error: %w", err)
				}

				_result = _r // avoid compiler optimizations
			}
		})
	}
}

func TestCalcDeps(t *testing.T) {

	tests := map[string]struct {
		tasks    map[string]*v1alpha1.PipelineTask
		assertFn func(*testing.T, map[string]*v1alpha1.PipelineTask, error)
	}{
		"Simple1": {
			tasks: map[string]*v1alpha1.PipelineTask{
				"task1": newTaskWithNeeds("task1", ""),
				"task2": newTaskWithNeeds("task2", "task1"),
			},
			assertFn: func(t *testing.T, tasks map[string]*v1alpha1.PipelineTask, _ error) {
				assert.Len(t, tasks["task1"].Dependants, 1)
				assert.Equal(t, tasks["task1"].Dependants[0].Name, "task2")
			},
		},
		"Error Circular": {
			tasks: map[string]*v1alpha1.PipelineTask{
				"task1": newTaskWithNeeds("task1", "task1"),
			},
			assertFn: func(t *testing.T, tasks map[string]*v1alpha1.PipelineTask, err error) {
				assert.Error(t, err, "circular dependency in task: task1")
			},
		},
		"Complex1": {
			tasks: map[string]*v1alpha1.PipelineTask{
				"task1": newTaskWithNeeds("task1", ""),
				"task2": newTaskWithNeeds("task2", "task1"),
				"task3": newTaskWithNeeds("task3", "task1"),
				"task4": newTaskWithNeeds("task4", "task1 && task2 && !task3"),
			},
			assertFn: func(t *testing.T, tasks map[string]*v1alpha1.PipelineTask, _ error) {
				assert.Len(t, tasks["task1"].Dependants, 3)
				assert.Contains(t, tasks["task1"].Dependants, tasks["task2"])
				assert.Contains(t, tasks["task1"].Dependants, tasks["task3"])
				assert.Contains(t, tasks["task1"].Dependants, tasks["task4"])

				assert.Len(t, tasks["task2"].Dependants, 1)
				assert.Contains(t, tasks["task2"].Dependants, tasks["task4"])

				assert.Len(t, tasks["task3"].Dependants, 1)
				assert.Contains(t, tasks["task3"].Dependants, tasks["task4"])

				assert.Len(t, tasks["task4"].Dependants, 0)
			},
		},
	}
	for tname, tt := range tests {
		t.Run(tname, func(t *testing.T) {
			tt.assertFn(t, tt.tasks, CalcDeps(tt.tasks))
		})
	}
}

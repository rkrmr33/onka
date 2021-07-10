package bexpr

import (
	"github.com/rkrmr33/onka/jag/compiler/bexpr/grammar"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
)

//go:generate pigeon -o grammar/grammar.go -optimize-parser grammar/grammar.peg
//go:generate goimports -w grammar/grammar.go

// EvalNeedsExpr takes a 'needs' expression and a map of all the task statuses, parses the
// expression and returns 'true' if the task with the 'needs' expression should be executed
// and 'false' or non-nil error otherwise.
func EvalNeedsExpr(t *v1alpha1.PipelineTask, statuses map[string]*v1alpha1.TaskStatus, opts ...grammar.Option) (bool, error) {
	if t.Task.Needs == "" {
		return true, nil
	}

	opts = append(opts,
		grammar.GlobalStore("statuses", statuses),
		grammar.GlobalStore("current", t.Name),
	)

	res, err := grammar.Parse("", []byte(t.Task.Needs), opts...)
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}

// CalcDeps calculates the dependencies between all tasks, setting the correct Dependent tasks
// pointers in each corresponding task
func CalcDeps(tasks map[string]*v1alpha1.PipelineTask) error {
	statuses := make(map[string]*v1alpha1.TaskStatus, len(tasks))
	for _, t := range tasks {
		statuses[t.Name] = &v1alpha1.TaskStatus{
			State: v1alpha1.TaskState_TASK_STATE_UNSPECIFIED,
		}
	}

	for _, t := range tasks {
		_, err := EvalNeedsExpr(t, statuses, grammar.GlobalStore("tasks", tasks))
		if err != nil {
			return err
		}
	}

	return nil
}

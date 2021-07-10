package compiler

import (
	"github.com/rkrmr33/onka/jag/compiler/bexpr"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
)

// calcEntrypoints receives all of the pipeline tasks and tries to build a tree
// of steps, considering the order they should execute in.
func calcEntrypoints(tasks []v1alpha1.PipelineTask) ([]v1alpha1.PipelineTask, error) {
	if len(tasks) == 0 {
		return nil, ErrNoEntry
	}

	if isPipelineSync(tasks) {
		return []v1alpha1.PipelineTask{buildSyncTree(tasks)}, nil
	}

	return buildParallelTree(tasks)
}

func isPipelineSync(tasks []v1alpha1.PipelineTask) bool {
	for _, t := range tasks {
		if t.Task.Needs != "" {
			return false
		}
	}
	return true
}

func buildSyncTree(tasks []v1alpha1.PipelineTask) v1alpha1.PipelineTask {
	var (
		cur  *v1alpha1.PipelineTask
		prev *v1alpha1.PipelineTask
	)

	head := v1alpha1.PipelineTask{
		Name: tasks[0].Name,
		Task: tasks[0].Task,
	}
	prev = &head

	for i := 1; i < len(tasks); i++ {
		cur = &v1alpha1.PipelineTask{
			Name: tasks[i].Name,
			Task: tasks[i].Task,
		}
		cur.Task.Needs = prev.Name
		prev.Dependants = []*v1alpha1.PipelineTask{cur}
		prev = cur
	}

	return head
}

func buildParallelTree(tasks []v1alpha1.PipelineTask) ([]v1alpha1.PipelineTask, error) {
	tasksMap := make(map[string]*v1alpha1.PipelineTask, len(tasks))
	for i := range tasks {
		tasksMap[tasks[i].Name] = &tasks[i]
	}

	if err := bexpr.CalcDeps(tasksMap); err != nil {
		return nil, err
	}

	entrypoints := findEntrypoints(tasks)
	if len(entrypoints) == 0 {
		return nil, ErrNoEntry
	}

	return entrypoints, nil
}

func findEntrypoints(tasks []v1alpha1.PipelineTask) []v1alpha1.PipelineTask {
	res := make([]v1alpha1.PipelineTask, 0)
	for _, t := range tasks {
		if t.Task.Needs == "" {
			res = append(res, t)
		}
	}
	return res
}

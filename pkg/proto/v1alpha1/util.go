package v1alpha1

import (
	"gopkg.in/yaml.v3"
)

type PipelineTask struct {
	Name string
	Task *TaskSpec

	Dependants []*PipelineTask
}

func (p *Pipeline) GetOrederdTasksYAML(data []byte) ([]PipelineTask, error) {
	root := &yaml.Node{}
	if err := yaml.Unmarshal(data, root); err != nil {
		return nil, err
	}

	ret := make([]PipelineTask, len(p.Spec.Tasks))
	tasksNode := getNodeWithName(root, "tasks")

	for i := 0; i < len(tasksNode.Content); i += 2 {
		nameNode := tasksNode.Content[i]

		ret[i/2] = PipelineTask{
			Name:       nameNode.Value,
			Task:       p.Spec.Tasks[nameNode.Value],
			Dependants: make([]*PipelineTask, 0),
		}
	}

	return ret, nil
}

func getNodeWithName(root *yaml.Node, name string) *yaml.Node {
	for i, n := range root.Content {
		if n.Tag == "!!str" && n.Value == name {
			return root.Content[i+1]
		}
		ret := getNodeWithName(n, name)
		if ret != nil {
			return ret
		}
	}

	return nil
}

func (s TaskState) IsFinal() bool {
	switch s {
	case TaskState_TASK_STATE_FAILURE, TaskState_TASK_STATE_SUCCESS:
		return true
	}
	return false
}

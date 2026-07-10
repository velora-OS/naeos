package scheduler

import (
	"fmt"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
)

type Scheduler interface {
	Schedule(neir any) ([]Task, error)
}

type Task struct {
	ID           string
	Name         string
	Dependencies []string
	Priority     int
}

type DefaultScheduler struct{}

func NewScheduler() Scheduler {
	return DefaultScheduler{}
}

func (DefaultScheduler) Schedule(neir any) ([]Task, error) {
	if neir == nil {
		return nil, fmt.Errorf("neir is nil")
	}

	neirModel, ok := neir.(*model.NEIR)
	if !ok {
		return []Task{{ID: "task-1", Name: "bootstrap"}}, nil
	}

	var tasks []Task
	taskID := 0

	tasks = append(tasks, Task{
		ID:       fmt.Sprintf("task-%d", taskID),
		Name:     "validate-specification",
		Priority: 0,
	})
	taskID++

	tasks = append(tasks, Task{
		ID:           fmt.Sprintf("task-%d", taskID),
		Name:         "build-model",
		Dependencies: []string{fmt.Sprintf("task-%d", taskID-1)},
		Priority:     1,
	})
	taskID++

	for _, mod := range neirModel.Modules {
		dep := fmt.Sprintf("task-%d", taskID-1)
		tasks = append(tasks, Task{
			ID:           fmt.Sprintf("task-%d", taskID),
			Name:         fmt.Sprintf("generate-module-%s", mod.Name),
			Dependencies: []string{dep},
			Priority:     2,
		})
		taskID++
	}

	if len(neirModel.Modules) > 0 {
		lastModuleTask := fmt.Sprintf("task-%d", taskID-1)
		for _, svc := range neirModel.Services {
			tasks = append(tasks, Task{
				ID:           fmt.Sprintf("task-%d", taskID),
				Name:         fmt.Sprintf("generate-service-%s", svc.Name),
				Dependencies: []string{lastModuleTask},
				Priority:     3,
			})
			taskID++
		}
	}

	if len(neirModel.Services) > 0 {
		lastServiceTask := fmt.Sprintf("task-%d", taskID-1)
		tasks = append(tasks, Task{
			ID:           fmt.Sprintf("task-%d", taskID),
			Name:         "generate-config",
			Dependencies: []string{lastServiceTask},
			Priority:     4,
		})
		taskID++
	} else if len(neirModel.Modules) > 0 {
		lastModuleTask := fmt.Sprintf("task-%d", taskID-1)
		tasks = append(tasks, Task{
			ID:           fmt.Sprintf("task-%d", taskID),
			Name:         "generate-config",
			Dependencies: []string{lastModuleTask},
			Priority:     4,
		})
		taskID++
	}

	tasks = append(tasks, Task{
		ID:       fmt.Sprintf("task-%d", taskID),
		Name:     "validate-output",
		Priority: 5,
	})

	return tasks, nil
}

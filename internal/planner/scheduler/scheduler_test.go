package scheduler

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/service"
)

func TestSchedulerReturnsErrorForNilInput(t *testing.T) {
	s := NewScheduler()
	_, err := s.Schedule(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestSchedulerGeneratesTasksFromNEIR(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "acme-api"},
		Modules: []module.Module{
			{Name: "auth", Path: "./internal/auth"},
			{Name: "user", Path: "./internal/user"},
		},
		Services: []service.Service{
			{Name: "api", Kind: "http", Port: 8080},
		},
	}

	s := NewScheduler()
	tasks, err := s.Schedule(neir)
	if err != nil {
		t.Fatalf("Schedule returned error: %v", err)
	}

	if len(tasks) < 5 {
		t.Fatalf("expected at least 5 tasks, got %d", len(tasks))
	}

	if tasks[0].Name != "validate-specification" {
		t.Fatalf("expected first task validate-specification, got %s", tasks[0].Name)
	}

	if tasks[1].Name != "build-model" {
		t.Fatalf("expected second task build-model, got %s", tasks[1].Name)
	}

	foundAuth := false
	foundUser := false
	for _, task := range tasks {
		if task.Name == "generate-module-auth" {
			foundAuth = true
		}
		if task.Name == "generate-module-user" {
			foundUser = true
		}
	}
	if !foundAuth || !foundUser {
		t.Fatalf("expected module tasks for auth and user, got tasks: %v", tasks)
	}

	lastTask := tasks[len(tasks)-1]
	if lastTask.Name != "validate-output" {
		t.Fatalf("expected last task validate-output, got %s", lastTask.Name)
	}
}

func TestSchedulerWithoutServices(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "simple-app"},
		Modules: []module.Module{
			{Name: "core", Path: "./internal/core"},
		},
	}

	s := NewScheduler()
	tasks, err := s.Schedule(neir)
	if err != nil {
		t.Fatalf("Schedule returned error: %v", err)
	}

	if len(tasks) < 4 {
		t.Fatalf("expected at least 4 tasks, got %d", len(tasks))
	}

	hasServiceTask := false
	for _, task := range tasks {
		if task.Name == "generate-service-api" {
			hasServiceTask = true
		}
	}
	if hasServiceTask {
		t.Fatal("should not have service task when no services defined")
	}
}

func TestSchedulerFallbackForNonNEIRInput(t *testing.T) {
	s := NewScheduler()
	tasks, err := s.Schedule("not-a-neir")
	if err != nil {
		t.Fatalf("Schedule returned error: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Name != "bootstrap" {
		t.Fatalf("expected fallback bootstrap task, got %v", tasks)
	}
}

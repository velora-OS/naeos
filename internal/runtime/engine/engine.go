package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Artifact struct {
	Path    string
	Content []byte
}

type ExecutionResult struct {
	Artifact Artifact
	Status   string
	Output   string
	Error    error
}

type RuntimeEngine interface {
	Run(artifact any) error
	Execute(artifact Artifact) (*ExecutionResult, error)
	ExecuteAll(artifacts []Artifact) ([]ExecutionResult, error)
	Validate(artifact Artifact) error
	SetOutputDir(dir string)
}

type DefaultRuntimeEngine struct {
	mu        sync.Mutex
	history   []ExecutionResult
	executed  map[string]bool
	outputDir string
}

func NewEngine() RuntimeEngine {
	return &DefaultRuntimeEngine{
		executed: make(map[string]bool),
	}
}

func (e *DefaultRuntimeEngine) SetOutputDir(dir string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.outputDir = dir
}

func (e *DefaultRuntimeEngine) Run(artifact any) error {
	if artifact == nil {
		return fmt.Errorf("artifact is nil")
	}
	return nil
}

func (e *DefaultRuntimeEngine) Execute(artifact Artifact) (*ExecutionResult, error) {
	if artifact.Path == "" {
		return nil, fmt.Errorf("artifact path must not be empty")
	}

	if err := e.Validate(artifact); err != nil {
		return nil, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	result := &ExecutionResult{
		Artifact: artifact,
		Status:   "completed",
	}

	if e.executed[artifact.Path] {
		result.Status = "skipped"
		result.Output = "already executed"
		e.history = append(e.history, *result)
		return result, nil
	}

	if e.outputDir != "" {
		fullPath := filepath.Join(e.outputDir, artifact.Path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			result.Status = "failed"
			result.Error = fmt.Errorf("create directory %s: %w", dir, err)
			e.history = append(e.history, *result)
			return result, result.Error
		}
		if err := os.WriteFile(fullPath, artifact.Content, 0o600); err != nil {
			result.Status = "failed"
			result.Error = fmt.Errorf("write file %s: %w", fullPath, err)
			e.history = append(e.history, *result)
			return result, result.Error
		}
		result.Output = fmt.Sprintf("wrote %s (%d bytes)", fullPath, len(artifact.Content))
	} else {
		result.Output = fmt.Sprintf("executed %s (%d bytes)", artifact.Path, len(artifact.Content))
	}

	e.executed[artifact.Path] = true
	e.history = append(e.history, *result)
	return result, nil
}

func (e *DefaultRuntimeEngine) ExecuteAll(artifacts []Artifact) ([]ExecutionResult, error) {
	if len(artifacts) == 0 {
		return nil, fmt.Errorf("no artifacts to execute")
	}

	var results []ExecutionResult
	for _, artifact := range artifacts {
		result, err := e.Execute(artifact)
		if err != nil {
			return results, fmt.Errorf("failed to execute %s: %w", artifact.Path, err)
		}
		results = append(results, *result)
	}
	return results, nil
}

func (e *DefaultRuntimeEngine) Validate(artifact Artifact) error {
	if artifact.Path == "" {
		return fmt.Errorf("artifact path must not be empty")
	}

	ext := filepath.Ext(artifact.Path)
	switch ext {
	case ".go":
		if len(artifact.Content) == 0 {
			return fmt.Errorf("go file %s has no content", artifact.Path)
		}
		content := string(artifact.Content)
		if !strings.Contains(content, "package ") {
			return fmt.Errorf("go file %s missing package declaration", artifact.Path)
		}
	case ".yaml", ".yml":
		if len(artifact.Content) == 0 {
			return fmt.Errorf("yaml file %s has no content", artifact.Path)
		}
	case ".md":
		if len(artifact.Content) == 0 {
			return fmt.Errorf("markdown file %s has no content", artifact.Path)
		}
	}

	return nil
}

func (e *DefaultRuntimeEngine) History() []ExecutionResult {
	e.mu.Lock()
	defer e.mu.Unlock()

	result := make([]ExecutionResult, len(e.history))
	copy(result, e.history)
	return result
}

func (e *DefaultRuntimeEngine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.history = nil
	e.executed = make(map[string]bool)
}

func (e *DefaultRuntimeEngine) ExecutedCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.executed)
}

func (e *DefaultRuntimeEngine) FailedCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	count := 0
	for _, r := range e.history {
		if r.Status == "failed" {
			count++
		}
	}
	return count
}

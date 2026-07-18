package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestWorkspace(t *testing.T) (string, *WorkspaceManager) {
	t.Helper()
	dir := t.TempDir()
	m := NewManager(dir)
	return dir, m
}

func TestInitAndLoad(t *testing.T) {
	dir, m := setupTestWorkspace(t)
	ws, err := m.Init("test-workspace")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if ws.Name != "test-workspace" {
		t.Errorf("expected name test-workspace, got %s", ws.Name)
	}
	if ws.Root != dir {
		t.Errorf("expected root %s, got %s", dir, ws.Root)
	}
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != ws.Name {
		t.Errorf("loaded name mismatch: %s vs %s", loaded.Name, ws.Name)
	}
}

func TestLoadNonexistent(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)
	_, err := m.Load()
	if err == nil {
		t.Fatal("expected error loading nonexistent workspace")
	}
}

func TestAddAndRemoveModule(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "mod1"), []byte("code"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "mod2"), []byte("code"), 0o644)

	if err := m.AddModule("mod1", "mod1", "", nil); err != nil {
		t.Fatalf("AddModule: %v", err)
	}
	if err := m.AddModule("mod2", "mod2", "", []string{"mod1"}); err != nil {
		t.Fatalf("AddModule with dep: %v", err)
	}
	if err := m.AddModule("mod1", "mod1", "", nil); err == nil {
		t.Fatal("expected error adding duplicate module")
	}
	if err := m.AddModule("mod3", "mod3", "", []string{"nonexistent"}); err == nil {
		t.Fatal("expected error adding module with nonexistent dep")
	}
	modules, _ := m.ListModules()
	if len(modules) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(modules))
	}
	if err := m.RemoveModule("mod1"); err == nil {
		t.Fatal("expected error removing module with dependents")
	}
	if err := m.RemoveModule("mod2"); err != nil {
		t.Fatalf("RemoveModule: %v", err)
	}
	if err := m.RemoveModule("nonexistent"); err == nil {
		t.Fatal("expected error removing nonexistent module")
	}
}

func TestGetModule(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "mod1"), []byte("code"), 0o644)
	m.AddModule("mod1", "mod1", "", nil)

	mod, err := m.GetModule("mod1")
	if err != nil {
		t.Fatalf("GetModule: %v", err)
	}
	if mod.Name != "mod1" {
		t.Errorf("expected mod1, got %s", mod.Name)
	}
	_, err = m.GetModule("nonexistent")
	if err == nil {
		t.Fatal("expected error getting nonexistent module")
	}
}

func TestDependencyGraph(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "a"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "b"), []byte("b"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "c"), []byte("c"), 0o644)
	m.AddModule("a", "a", "", nil)
	m.AddModule("b", "b", "", []string{"a"})
	m.AddModule("c", "c", "", []string{"b"})

	g, err := m.DependencyGraph()
	if err != nil {
		t.Fatalf("DependencyGraph: %v", err)
	}
	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort: %v", err)
	}
	if len(order) != 3 {
		t.Fatalf("expected 3 in order, got %d", len(order))
	}
	aIdx, bIdx, cIdx := -1, -1, -1
	for i, name := range order {
		switch name {
		case "a":
			aIdx = i
		case "b":
			bIdx = i
		case "c":
			cIdx = i
		}
	}
	if aIdx >= bIdx || bIdx >= cIdx {
		t.Errorf("topological order wrong: a=%d, b=%d, c=%d", aIdx, bIdx, cIdx)
	}
	deps := g.TransitiveDeps("c")
	if len(deps) != 2 {
		t.Errorf("expected 2 transitive deps for c, got %d", len(deps))
	}
	dependents := g.Dependents("a")
	if len(dependents) != 1 || dependents[0] != "b" {
		t.Errorf("expected [b] dependents for a, got %v", dependents)
	}
}

func TestCircularDependency(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "a"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "b"), []byte("b"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "c"), []byte("c"), 0o644)
	m.AddModule("a", "a", "", nil)
	m.AddModule("b", "b", "", nil)
	m.AddModule("c", "c", "", nil)

	ws, _ := m.Load()
	ws.Modules[0].DependsOn = []string{"c"}
	ws.Modules[1].DependsOn = []string{"a"}
	ws.Modules[2].DependsOn = []string{"b"}
	data, _ := json.MarshalIndent(ws, "", "  ")
	os.WriteFile(m.configPath(), data, 0o600)

	g, _ := m.DependencyGraph()
	_, err := g.TopologicalSort()
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
	cycles := g.CircularDependencies()
	if len(cycles) == 0 {
		t.Fatal("expected circular dependencies detected")
	}
}

func TestLockFile(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "mod1"), []byte("content1"), 0o644)
	m.AddModule("mod1", "mod1", "", nil)

	if err := m.SaveLock(); err != nil {
		t.Fatalf("SaveLock: %v", err)
	}
	lock, err := m.LoadLock()
	if err != nil {
		t.Fatalf("LoadLock: %v", err)
	}
	if lock.WorkspaceName != "test" {
		t.Errorf("expected workspace name test, got %s", lock.WorkspaceName)
	}
	info, ok := lock.Modules["mod1"]
	if !ok {
		t.Fatal("expected mod1 in lock")
	}
	if info.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestDetectDirty(t *testing.T) {
	dir, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(dir, "mod1"), []byte("content1"), 0o644)
	m.AddModule("mod1", "mod1", "", nil)

	dirty, err := m.DetectDirty()
	if err != nil {
		t.Fatalf("DetectDirty without lock: %v", err)
	}
	if len(dirty) != 1 || !dirty[0].IsNew {
		t.Error("expected mod1 to be new")
	}

	m.SaveLock()
	dirty, err = m.DetectDirty()
	if err != nil {
		t.Fatalf("DetectDirty: %v", err)
	}
	if len(dirty) != 0 {
		t.Errorf("expected no dirty modules, got %d", len(dirty))
	}

	os.WriteFile(filepath.Join(dir, "mod1"), []byte("changed"), 0o644)
	dirty, err = m.DetectDirty()
	if err != nil {
		t.Fatalf("DetectDirty after change: %v", err)
	}
	if len(dirty) != 1 || dirty[0].IsNew {
		t.Error("expected mod1 to be dirty (modified)")
	}
}

func TestValidate(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "a"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "b"), []byte("b"), 0o644)
	m.AddModule("a", "a", "", nil)
	m.AddModule("b", "b", "", []string{"a"})

	issues, err := m.Validate()
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	errCount := 0
	for _, issue := range issues {
		if issue.Severity == "error" {
			errCount++
		}
	}
	if errCount != 0 {
		t.Errorf("expected 0 errors, got %d: %v", errCount, issues)
	}

	ws, _ := m.Load()
	ws.Modules = append(ws.Modules, ModuleRef{Name: "a", Path: "a"})
	data, _ := json.MarshalIndent(ws, "", "  ")
	os.WriteFile(m.configPath(), data, 0o600)

	issues, _ = m.Validate()
	hasDup := false
	for _, issue := range issues {
		if issue.Severity == "error" {
			hasDup = true
		}
	}
	if !hasDup {
		t.Error("expected duplicate module error")
	}
}

func TestReport(t *testing.T) {
	_, m := setupTestWorkspace(t)
	m.Init("test")
	os.WriteFile(filepath.Join(m.rootDir, "a"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(m.rootDir, "b"), []byte("b"), 0o644)
	m.AddModule("a", "a", "", nil)
	m.AddModule("b", "b", "", []string{"a"})

	report, err := m.Report()
	if err != nil {
		t.Fatalf("Report: %v", err)
	}
	if report.ModuleCount != 2 {
		t.Errorf("expected 2 modules, got %d", report.ModuleCount)
	}
	if len(report.BuildOrder) != 2 {
		t.Errorf("expected 2 in build order, got %d", len(report.BuildOrder))
	}
	if report.HasCircular {
		t.Error("expected no circular dependencies")
	}
	if len(report.Modules) != 2 {
		t.Errorf("expected 2 module reports, got %d", len(report.Modules))
	}
}

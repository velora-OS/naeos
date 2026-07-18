package workspace

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Workspace struct {
	Name        string         `json:"name"`
	Root        string         `json:"root"`
	Modules     []ModuleRef    `json:"modules"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Description string         `json:"description,omitempty"`
}

type ModuleRef struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	SpecFile   string   `json:"spec_file,omitempty"`
	DependsOn  []string `json:"depends_on,omitempty"`
	Version    string   `json:"version,omitempty"`
	Checksum   string   `json:"checksum,omitempty"`
}

type WorkspaceManager struct {
	rootDir string
	mu      sync.RWMutex
}

func NewManager(rootDir string) *WorkspaceManager {
	return &WorkspaceManager{rootDir: rootDir}
}

func (m *WorkspaceManager) configPath() string {
	return filepath.Join(m.rootDir, "naeos.workspace.json")
}

func (m *WorkspaceManager) lockPath() string {
	return filepath.Join(m.rootDir, "naeos.workspace.lock")
}

func (m *WorkspaceManager) Init(name string) (*Workspace, error) {
	ws := &Workspace{
		Name:      name,
		Root:      m.rootDir,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(m.rootDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(m.configPath(), data, 0o600); err != nil {
		return nil, err
	}
	return ws, nil
}

func (m *WorkspaceManager) Load() (*Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, err := os.ReadFile(m.configPath())
	if err != nil {
		return nil, fmt.Errorf("no workspace found at %s: %w", m.rootDir, err)
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (m *WorkspaceManager) AddModule(name, path, specFile string, dependsOn []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ws, err := m.loadUnsafe()
	if err != nil {
		return err
	}
	for _, mod := range ws.Modules {
		if mod.Name == name {
			return fmt.Errorf("module %s already exists", name)
		}
	}
	for _, dep := range dependsOn {
		found := false
		for _, existing := range ws.Modules {
			if existing.Name == dep {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("dependency %s not found in workspace", dep)
		}
	}
	checksum, err := m.computeModuleChecksum(path)
	if err != nil {
		checksum = ""
	}
	ws.Modules = append(ws.Modules, ModuleRef{
		Name:      name,
		Path:      path,
		SpecFile:  specFile,
		DependsOn: dependsOn,
		Checksum:  checksum,
	})
	ws.UpdatedAt = time.Now()
	return m.saveUnsafe(ws)
}

func (m *WorkspaceManager) RemoveModule(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ws, err := m.loadUnsafe()
	if err != nil {
		return err
	}
	for _, mod := range ws.Modules {
		if mod.Name != name {
			for _, dep := range mod.DependsOn {
				if dep == name {
					return fmt.Errorf("cannot remove %s: module %s depends on it", name, mod.Name)
				}
			}
		}
	}
	for i, mod := range ws.Modules {
		if mod.Name == name {
			ws.Modules = append(ws.Modules[:i], ws.Modules[i+1:]...)
			ws.UpdatedAt = time.Now()
			return m.saveUnsafe(ws)
		}
	}
	return fmt.Errorf("module %s not found", name)
}

func (m *WorkspaceManager) ListModules() ([]ModuleRef, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	return ws.Modules, nil
}

func (m *WorkspaceManager) GetModule(name string) (*ModuleRef, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	for _, mod := range ws.Modules {
		if mod.Name == name {
			return &mod, nil
		}
	}
	return nil, fmt.Errorf("module %s not found", name)
}

func (m *WorkspaceManager) loadUnsafe() (*Workspace, error) {
	data, err := os.ReadFile(m.configPath())
	if err != nil {
		return nil, fmt.Errorf("no workspace found at %s: %w", m.rootDir, err)
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, err
	}
	return &ws, nil
}

func (m *WorkspaceManager) saveUnsafe(ws *Workspace) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath(), data, 0o600)
}

func (m *WorkspaceManager) computeModuleChecksum(path string) (string, error) {
	fullPath := filepath.Join(m.rootDir, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return m.computeDirChecksum(fullPath)
	}
	return m.computeFileChecksum(fullPath)
}

func (m *WorkspaceManager) computeFileChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:8]), nil
}

func (m *WorkspaceManager) computeDirChecksum(dir string) (string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	h := sha256.New()
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		h.Write(data)
	}
	return fmt.Sprintf("%x", h.Sum(nil)[:8]), nil
}

type DependencyGraph struct {
	Edges    map[string][]string
	Vertices map[string]*ModuleRef
}

func (m *WorkspaceManager) DependencyGraph() (*DependencyGraph, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	g := &DependencyGraph{
		Edges:    make(map[string][]string),
		Vertices: make(map[string]*ModuleRef),
	}
	for i := range ws.Modules {
		g.Vertices[ws.Modules[i].Name] = &ws.Modules[i]
		g.Edges[ws.Modules[i].Name] = ws.Modules[i].DependsOn
	}
	return g, nil
}

func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var order []string
	var visit func(name string) error
	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving %s", name)
		}
		visiting[name] = true
		for _, dep := range g.Edges[name] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visiting[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}
	for name := range g.Vertices {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (g *DependencyGraph) CircularDependencies() [][]string {
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var cycles [][]string
	var path []string
	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		if visiting[name] {
			cycleStart := -1
			for i, p := range path {
				if p == name {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				cycles = append(cycles, cycle)
			}
			return
		}
		visiting[name] = true
		path = append(path, name)
		for _, dep := range g.Edges[name] {
			visit(dep)
		}
		path = path[:len(path)-1]
		visiting[name] = false
		visited[name] = true
	}
	for name := range g.Vertices {
		visit(name)
	}
	return cycles
}

func (g *DependencyGraph) Dependents(module string) []string {
	var dependents []string
	for name, deps := range g.Edges {
		for _, dep := range deps {
			if dep == module {
				dependents = append(dependents, name)
				break
			}
		}
	}
	sort.Strings(dependents)
	return dependents
}

func (g *DependencyGraph) TransitiveDeps(module string) []string {
	visited := make(map[string]bool)
	var deps []string
	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true
		for _, dep := range g.Edges[name] {
			deps = append(deps, dep)
			visit(dep)
		}
	}
	visit(module)
	sort.Strings(deps)
	return deps
}

type LockFile struct {
	WorkspaceName string              `json:"workspace_name"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Modules       map[string]LockInfo `json:"modules"`
}

type LockInfo struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
	Version  string `json:"version,omitempty"`
}

func (m *WorkspaceManager) SaveLock() error {
	ws, err := m.Load()
	if err != nil {
		return err
	}
	lock := &LockFile{
		WorkspaceName: ws.Name,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Modules:       make(map[string]LockInfo),
	}
	for _, mod := range ws.Modules {
		checksum, err := m.computeModuleChecksum(mod.Path)
		if err != nil {
			checksum = ""
		}
		lock.Modules[mod.Name] = LockInfo{
			Path:     mod.Path,
			Checksum: checksum,
			Version:  mod.Version,
		}
	}
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.lockPath(), data, 0o600)
}

func (m *WorkspaceManager) LoadLock() (*LockFile, error) {
	data, err := os.ReadFile(m.lockPath())
	if err != nil {
		return nil, fmt.Errorf("no lock file found: %w", err)
	}
	var lock LockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, err
	}
	return &lock, nil
}

type DirtyInfo struct {
	Module     string `json:"module"`
	OldSum    string `json:"old_sum"`
	NewSum    string `json:"new_sum"`
	IsNew     bool   `json:"is_new"`
	IsRemoved bool   `json:"is_removed"`
}

func (m *WorkspaceManager) DetectDirty() ([]DirtyInfo, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	lock, err := m.LoadLock()
	if err != nil {
		var dirty []DirtyInfo
		for _, mod := range ws.Modules {
			checksum, err := m.computeModuleChecksum(mod.Path)
			if err != nil {
				checksum = ""
			}
			dirty = append(dirty, DirtyInfo{
				Module: mod.Name,
				NewSum: checksum,
				IsNew:  true,
			})
		}
		return dirty, nil
	}
	var dirty []DirtyInfo
	currentModules := make(map[string]bool)
	for _, mod := range ws.Modules {
		currentModules[mod.Name] = true
		checksum, err := m.computeModuleChecksum(mod.Path)
		if err != nil {
			checksum = ""
		}
		lockInfo, exists := lock.Modules[mod.Name]
		if !exists {
			dirty = append(dirty, DirtyInfo{
				Module: mod.Name,
				NewSum: checksum,
				IsNew:  true,
			})
		} else if lockInfo.Checksum != checksum {
			dirty = append(dirty, DirtyInfo{
				Module: mod.Name,
				OldSum: lockInfo.Checksum,
				NewSum: checksum,
			})
		}
	}
	for name := range lock.Modules {
		if !currentModules[name] {
			dirty = append(dirty, DirtyInfo{
				Module:     name,
				OldSum:     lock.Modules[name].Checksum,
				IsRemoved: true,
			})
		}
	}
	return dirty, nil
}

type ValidationIssue struct {
	Module   string `json:"module"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func (m *WorkspaceManager) Validate() ([]ValidationIssue, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	var issues []ValidationIssue
	names := make(map[string]bool)
	for _, mod := range ws.Modules {
		if names[mod.Name] {
			issues = append(issues, ValidationIssue{
				Module:   mod.Name,
				Severity: "error",
				Message:  fmt.Sprintf("duplicate module name: %s", mod.Name),
			})
		}
		names[mod.Name] = true
		fullPath := filepath.Join(m.rootDir, mod.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			issues = append(issues, ValidationIssue{
				Module:   mod.Name,
				Severity: "error",
				Message:  fmt.Sprintf("path does not exist: %s", mod.Path),
			})
		}
		if mod.SpecFile != "" {
			specPath := filepath.Join(m.rootDir, mod.SpecFile)
			if _, err := os.Stat(specPath); os.IsNotExist(err) {
				issues = append(issues, ValidationIssue{
					Module:   mod.Name,
					Severity: "warning",
					Message:  fmt.Sprintf("spec file does not exist: %s", mod.SpecFile),
				})
			}
		}
		for _, dep := range mod.DependsOn {
			if !names[dep] {
				issues = append(issues, ValidationIssue{
					Module:   mod.Name,
					Severity: "error",
					Message:  fmt.Sprintf("dependency %s not found in workspace", dep),
				})
			}
		}
	}
	g, err := m.DependencyGraph()
	if err == nil {
		cycles := g.CircularDependencies()
		for _, cycle := range cycles {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("circular dependency: %s", strings.Join(cycle, " -> ")),
			})
		}
	}
	return issues, nil
}

type WorkspaceReport struct {
	Name          string         `json:"name"`
	Root          string         `json:"root"`
	ModuleCount   int            `json:"module_count"`
	BuildOrder    []string       `json:"build_order"`
	Modules       []ModuleReport `json:"modules"`
	HasCircular   bool           `json:"has_circular"`
	GeneratedAt   time.Time      `json:"generated_at"`
}

type ModuleReport struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	DependsOn  []string `json:"depends_on"`
	Dependents []string `json:"dependents"`
	Depth      int      `json:"depth"`
}

func (m *WorkspaceManager) Report() (*WorkspaceReport, error) {
	ws, err := m.Load()
	if err != nil {
		return nil, err
	}
	g, err := m.DependencyGraph()
	if err != nil {
		return nil, err
	}
	buildOrder, _ := g.TopologicalSort()
	cycles := g.CircularDependencies()
	report := &WorkspaceReport{
		Name:        ws.Name,
		Root:        ws.Root,
		ModuleCount: len(ws.Modules),
		BuildOrder:  buildOrder,
		HasCircular: len(cycles) > 0,
		GeneratedAt: time.Now(),
	}
	for _, mod := range ws.Modules {
		dependents := g.Dependents(mod.Name)
		deps := g.TransitiveDeps(mod.Name)
		report.Modules = append(report.Modules, ModuleReport{
			Name:       mod.Name,
			Path:       mod.Path,
			DependsOn:  mod.DependsOn,
			Dependents: dependents,
			Depth:      len(deps),
		})
	}
	return report, nil
}

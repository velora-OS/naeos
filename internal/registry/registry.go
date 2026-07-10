package registry

import (
	"fmt"
	"sync"
)

type Entry struct {
	Name      string
	Version   string
	Category  string
	Component any
	Metadata  map[string]string
}

type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*Entry),
	}
}

func (r *Registry) Register(name string, component any) error {
	return r.RegisterWithMeta(name, "", "", component, nil)
}

func (r *Registry) RegisterWithMeta(name, version, category string, component any, metadata map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if name == "" {
		return fmt.Errorf("entry name must not be empty")
	}
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("entry %s already registered", name)
	}

	r.entries[name] = &Entry{
		Name:      name,
		Version:   version,
		Category:  category,
		Component: component,
		Metadata:  metadata,
	}
	return nil
}

func (r *Registry) Resolve(name string) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.entries[name]
	if !exists {
		return nil, fmt.Errorf("entry %s not found", name)
	}
	return entry.Component, nil
}

func (r *Registry) GetEntry(name string) (*Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.entries[name]
	if !exists {
		return nil, fmt.Errorf("entry %s not found", name)
	}
	return entry, nil
}

func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[name]; !exists {
		return fmt.Errorf("entry %s not found", name)
	}
	delete(r.entries, name)
	return nil
}

func (r *Registry) RegisteredEntries() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]string, 0, len(r.entries))
	for name := range r.entries {
		result = append(result, name)
	}
	return result
}

func (r *Registry) FindByCategory(category string) []*Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Entry
	for _, entry := range r.entries {
		if entry.Category == category {
			result = append(result, entry)
		}
	}
	return result
}

func (r *Registry) FindByVersion(version string) []*Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Entry
	for _, entry := range r.entries {
		if entry.Version == version {
			result = append(result, entry)
		}
	}
	return result
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

func (r *Registry) Contains(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.entries[name]
	return exists
}

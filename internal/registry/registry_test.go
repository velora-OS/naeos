package registry

import (
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r.Count() != 0 {
		t.Fatalf("expected 0 entries, got %d", r.Count())
	}
}

func TestRegister(t *testing.T) {
	r := NewRegistry()
	err := r.Register("test-service", "service-component")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.Count())
	}
}

func TestRegisterEmptyName(t *testing.T) {
	r := NewRegistry()
	err := r.Register("", "component")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("test", "component")
	err := r.Register("test", "component")
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestResolve(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("test", "component")
	result, err := r.Resolve("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "component" {
		t.Fatalf("expected 'component', got %v", result)
	}
}

func TestResolveNotFound(t *testing.T) {
	r := NewRegistry()
	_, err := r.Resolve("missing")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestGetEntry(t *testing.T) {
	r := NewRegistry()
	_ = r.RegisterWithMeta("test", "1.0.0", "service", "component", map[string]string{"key": "value"})
	entry, err := r.GetEntry("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Name != "test" {
		t.Fatalf("expected name 'test', got %s", entry.Name)
	}
	if entry.Version != "1.0.0" {
		t.Fatalf("expected version '1.0.0', got %s", entry.Version)
	}
	if entry.Category != "service" {
		t.Fatalf("expected category 'service', got %s", entry.Category)
	}
	if entry.Metadata["key"] != "value" {
		t.Fatalf("expected metadata key=value, got %v", entry.Metadata)
	}
}

func TestUnregister(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("test", "component")
	err := r.Unregister("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Count() != 0 {
		t.Fatalf("expected 0 entries after unregister, got %d", r.Count())
	}
}

func TestUnregisterNotFound(t *testing.T) {
	r := NewRegistry()
	err := r.Unregister("missing")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestRegisteredEntries(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("a", "1")
	_ = r.Register("b", "2")
	_ = r.Register("c", "3")
	entries := r.RegisteredEntries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestFindByCategory(t *testing.T) {
	r := NewRegistry()
	_ = r.RegisterWithMeta("api", "1.0", "service", "api-component", nil)
	_ = r.RegisterWithMeta("db", "1.0", "storage", "db-component", nil)
	_ = r.RegisterWithMeta("web", "2.0", "service", "web-component", nil)

	services := r.FindByCategory("service")
	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(services))
	}
}

func TestFindByVersion(t *testing.T) {
	r := NewRegistry()
	_ = r.RegisterWithMeta("a", "1.0.0", "", "comp-a", nil)
	_ = r.RegisterWithMeta("b", "2.0.0", "", "comp-b", nil)
	_ = r.RegisterWithMeta("c", "1.0.0", "", "comp-c", nil)

	v1 := r.FindByVersion("1.0.0")
	if len(v1) != 2 {
		t.Fatalf("expected 2 entries with version 1.0.0, got %d", len(v1))
	}
}

func TestContains(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("test", "component")
	if !r.Contains("test") {
		t.Fatal("expected Contains to return true")
	}
	if r.Contains("missing") {
		t.Fatal("expected Contains to return false")
	}
}

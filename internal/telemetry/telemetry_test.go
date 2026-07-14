package telemetry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type mockExporter struct {
	spans [][]Span
	mu    sync.Mutex
}

func (m *mockExporter) ExportSpans(spans []Span) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	batch := make([]Span, len(spans))
	copy(batch, spans)
	m.spans = append(m.spans, batch)
	return nil
}

func (m *mockExporter) Flush() error { return nil }

func (m *mockExporter) totalSpans() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, batch := range m.spans {
		total += len(batch)
	}
	return total
}

func TestServiceStartEndSpan(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{BatchSize: 10}, mock)

	span := svc.StartSpan("test-span")
	if span.Name != "test-span" {
		t.Errorf("expected name 'test-span', got %s", span.Name)
	}
	if span.ID == "" {
		t.Error("expected non-empty ID")
	}

	svc.EndSpan(span)
	if svc.SpanCount() != 1 {
		t.Errorf("expected 1 buffered span, got %d", svc.SpanCount())
	}
}

func TestServiceAutoFlush(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{BatchSize: 2}, mock)

	for i := 0; i < 3; i++ {
		span := svc.StartSpan("span")
		svc.EndSpan(span)
	}

	if mock.totalSpans() < 2 {
		t.Errorf("expected at least 2 exported spans, got %d", mock.totalSpans())
	}
}

func TestServiceFlush(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{BatchSize: 100}, mock)

	span := svc.StartSpan("span")
	svc.EndSpan(span)

	if err := svc.Flush(); err != nil {
		t.Fatal(err)
	}
	if svc.SpanCount() != 0 {
		t.Errorf("expected 0 buffered spans after flush, got %d", svc.SpanCount())
	}
	if mock.totalSpans() != 1 {
		t.Errorf("expected 1 exported span, got %d", mock.totalSpans())
	}
}

func TestParentChildSpans(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{BatchSize: 100}, mock)

	parent := svc.StartSpan("parent")
	child := svc.StartSpanWithParent("child", parent.ID)

	if child.ParentID != parent.ID {
		t.Errorf("expected child ParentID=%s, got %s", parent.ID, child.ParentID)
	}
	if child.ID == parent.ID {
		t.Error("child and parent should have different IDs")
	}
}

func TestSpanLabels(t *testing.T) {
	span := &Span{
		Name:   "labeled",
		Labels: map[string]string{"env": "test", "service": "api"},
	}
	if span.Labels["env"] != "test" {
		t.Errorf("expected env=test, got %s", span.Labels["env"])
	}
}

func TestHTTPExporterNew(t *testing.T) {
	exp := NewHTTPExporter("http://localhost:9999", 0)
	if exp.client == nil {
		t.Fatal("expected non-nil client")
	}
	if exp.client.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", exp.client.Timeout)
	}
	if exp.endpoint != "http://localhost:9999" {
		t.Errorf("expected endpoint 'http://localhost:9999', got %s", exp.endpoint)
	}
}

func TestHTTPExporterFlushEmpty(t *testing.T) {
	exp := NewHTTPExporter("http://localhost:9999", time.Second)
	if err := exp.Flush(); err != nil {
		t.Fatalf("expected nil error on empty flush, got %v", err)
	}
}

func TestHTTPExporterExportSpans(t *testing.T) {
	var received []Span
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/traces" {
			t.Errorf("expected path /v1/traces, got %s", r.URL.Path)
		}
		var spans []Span
		if err := json.NewDecoder(r.Body).Decode(&spans); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		received = spans
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	exp := NewHTTPExporter(srv.URL, 5*time.Second)
	spans := []Span{
		{Name: "s1", ID: "id1", Status: "ok"},
		{Name: "s2", ID: "id2", Status: "ok"},
	}
	if err := exp.ExportSpans(spans); err != nil {
		t.Fatalf("ExportSpans failed: %v", err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 spans received by server, got %d", len(received))
	}
	if received[0].Name != "s1" || received[1].Name != "s2" {
		t.Errorf("unexpected span names: %v, %v", received[0].Name, received[1].Name)
	}
}

func TestHTTPExporterExportSpansError(t *testing.T) {
	exp := NewHTTPExporter("http://127.0.0.1:1", time.Second)
	spans := []Span{{Name: "s1", ID: "id1", Status: "ok"}}
	err := exp.ExportSpans(spans)
	if err == nil {
		t.Fatal("expected error for bad endpoint, got nil")
	}
}

func TestServiceDefaults(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{}, mock)
	if svc.config.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", svc.config.Timeout)
	}
	if svc.config.BatchSize != 100 {
		t.Errorf("expected default batch size 100, got %d", svc.config.BatchSize)
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()
	id3 := generateID()

	if id1 == id2 || id2 == id3 {
		t.Error("expected unique IDs from generateID")
	}
	// The counter is global and monotonically increasing.
	// Just verify the format and that they are ordered.
	if id1 >= id2 || id2 >= id3 {
		t.Errorf("expected strictly increasing IDs: %s < %s < %s", id1, id2, id3)
	}
}

func TestSpanCountAfterMultipleEndSpan(t *testing.T) {
	mock := &mockExporter{}
	svc := NewService(Config{BatchSize: 100}, mock)

	for i := 0; i < 5; i++ {
		span := svc.StartSpan("span")
		svc.EndSpan(span)
	}

	if svc.SpanCount() != 5 {
		t.Errorf("expected 5 buffered spans, got %d", svc.SpanCount())
	}
}

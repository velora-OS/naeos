package pluginsdk

import (
	"fmt"
	"sync"
	"time"
)

// Plugin Interface

type Plugin interface {
	Name() string
	Version() string
	Description() string
	Initialize(ctx *PluginContext) error
	Execute(action string, params map[string]interface{}) (interface{}, error)
	Shutdown() error
}

type PluginContext struct {
	Config    map[string]interface{}
	Logger    Logger
	Metrics   MetricsCollector
	EventBus  EventEmitter
}

type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type MetricsCollector interface {
	CounterInc(name string, labels map[string]string)
	GaugeSet(name string, value float64, labels map[string]string)
	HistogramObserve(name string, value float64, labels map[string]string)
}

type EventEmitter interface {
	Emit(event string, data interface{})
	On(event string, handler func(data interface{}))
}

// Plugin Manager

type Manager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

func (m *Manager) Register(plugin Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := plugin.Name()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}

	m.plugins[name] = plugin
	return nil
}

func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	delete(m.plugins, name)
	return nil
}

func (m *Manager) Get(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, ok := m.plugins[name]
	return plugin, ok
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

func (m *Manager) InitializeAll(ctx *PluginContext) error {
	for name, plugin := range m.plugins {
		if err := plugin.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize plugin '%s': %w", name, err)
		}
	}
	return nil
}

func (m *Manager) ShutdownAll() error {
	var lastErr error
	for name, plugin := range m.plugins {
		if err := plugin.Shutdown(); err != nil {
			lastErr = fmt.Errorf("failed to shutdown plugin '%s': %w", name, err)
		}
	}
	return lastErr
}

func (m *Manager) Execute(name, action string, params map[string]interface{}) (interface{}, error) {
	plugin, ok := m.Get(name)
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}
	return plugin.Execute(action, params)
}

// Plugin Manifest

type Manifest struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	License     string            `json:"license"`
	Dependencies []string         `json:"dependencies"`
	Actions     []ActionManifest `json:"actions"`
	Config      map[string]ConfigField `json:"config"`
}

type ActionManifest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Params      map[string]string `json:"params"`
	Returns     string            `json:"returns"`
}

type ConfigField struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
}

// Base Plugin (for embedding)

type BasePlugin struct {
	NameVal        string
	VersionVal     string
	DescriptionVal string
}

func (b *BasePlugin) Name() string        { return b.NameVal }
func (b *BasePlugin) Version() string     { return b.VersionVal }
func (b *BasePlugin) Description() string { return b.DescriptionVal }

func (b *BasePlugin) Initialize(ctx *PluginContext) error { return nil }
func (b *BasePlugin) Shutdown() error                     { return nil }

// Simple Logger

type SimpleLogger struct {
	prefix string
}

func NewSimpleLogger(prefix string) *SimpleLogger {
	return &SimpleLogger{prefix: prefix}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[%s] INFO: %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	fmt.Printf("[%s] WARN: %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[%s] ERROR: %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[%s] DEBUG: %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

// Simple Metrics

type SimpleMetrics struct{}

func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{}
}

func (m *SimpleMetrics) CounterInc(name string, labels map[string]string)    {}
func (m *SimpleMetrics) GaugeSet(name string, value float64, labels map[string]string) {}
func (m *SimpleMetrics) HistogramObserve(name string, value float64, labels map[string]string) {}

// Simple Event Emitter

type SimpleEventEmitter struct {
	handlers map[string][]func(data interface{})
	mu       sync.RWMutex
}

func NewSimpleEventEmitter() *SimpleEventEmitter {
	return &SimpleEventEmitter{
		handlers: make(map[string][]func(data interface{})),
	}
}

func (e *SimpleEventEmitter) Emit(event string, data interface{}) {
	e.mu.RLock()
	handlers := e.handlers[event]
	e.mu.RUnlock()

	for _, handler := range handlers {
		handler(data)
	}
}

func (e *SimpleEventEmitter) On(event string, handler func(data interface{})) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[event] = append(e.handlers[event], handler)
}

// Plugin Lifecycle

type PluginState string

const (
	StateCreated   PluginState = "created"
	StateInitialized PluginState = "initialized"
	StateRunning   PluginState = "running"
	StateStopped   PluginState = "stopped"
	StateError     PluginState = "error"
)

type PluginInfo struct {
	Name      string
	Version   string
	State     PluginState
	StartedAt time.Time
	Error     error
}

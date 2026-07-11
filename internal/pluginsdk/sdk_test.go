package pluginsdk

import (
	"testing"
)

type TestPlugin struct {
	BasePlugin
	initialized bool
	executed    string
}

func (p *TestPlugin) Execute(action string, params map[string]interface{}) (interface{}, error) {
	p.executed = action
	return "result", nil
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("expected manager to be created")
	}
}

func TestRegisterPlugin(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{
			NameVal:    "test",
			VersionVal: "1.0.0",
		},
	}

	err := m.Register(plugin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m.List()) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(m.List()))
	}
}

func TestRegisterDuplicatePlugin(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{NameVal: "test"},
	}

	m.Register(plugin)
	err := m.Register(plugin)
	if err == nil {
		t.Error("expected error for duplicate plugin")
	}
}

func TestUnregisterPlugin(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{NameVal: "test"},
	}

	m.Register(plugin)
	err := m.Unregister("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m.List()) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(m.List()))
	}
}

func TestUnregisterNotFound(t *testing.T) {
	m := NewManager()
	err := m.Unregister("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent plugin")
	}
}

func TestGetPlugin(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{NameVal: "test"},
	}

	m.Register(plugin)
	got, ok := m.Get("test")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if got.Name() != "test" {
		t.Errorf("expected name 'test', got %s", got.Name())
	}
}

func TestGetPluginNotFound(t *testing.T) {
	m := NewManager()
	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("expected plugin not found")
	}
}

func TestListPlugins(t *testing.T) {
	m := NewManager()
	m.Register(&TestPlugin{BasePlugin: BasePlugin{NameVal: "a"}})
	m.Register(&TestPlugin{BasePlugin: BasePlugin{NameVal: "b"}})

	list := m.List()
	if len(list) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(list))
	}
}

func TestExecutePlugin(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{NameVal: "test"},
	}

	m.Register(plugin)
	result, err := m.Execute("test", "doSomething", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "result" {
		t.Errorf("expected 'result', got %v", result)
	}

	if plugin.executed != "doSomething" {
		t.Errorf("expected 'doSomething', got %s", plugin.executed)
	}
}

func TestExecutePluginNotFound(t *testing.T) {
	m := NewManager()
	_, err := m.Execute("nonexistent", "action", nil)
	if err == nil {
		t.Error("expected error for nonexistent plugin")
	}
}

func TestInitializeAll(t *testing.T) {
	m := NewManager()
	plugin := &TestPlugin{
		BasePlugin: BasePlugin{NameVal: "test"},
	}

	m.Register(plugin)

	ctx := &PluginContext{
		Config:   map[string]interface{}{"key": "value"},
		Logger:   NewSimpleLogger("test"),
		Metrics:  NewSimpleMetrics(),
		EventBus: NewSimpleEventEmitter(),
	}

	err := m.InitializeAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShutdownAll(t *testing.T) {
	m := NewManager()
	m.Register(&TestPlugin{BasePlugin: BasePlugin{NameVal: "a"}})
	m.Register(&TestPlugin{BasePlugin: BasePlugin{NameVal: "b"}})

	err := m.ShutdownAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSimpleLogger(t *testing.T) {
	logger := NewSimpleLogger("test")
	logger.Info("test message %s", "arg")
	logger.Warn("test warning")
	logger.Error("test error")
	logger.Debug("test debug")
}

func TestSimpleEventEmitter(t *testing.T) {
	emitter := NewSimpleEventEmitter()

	received := false
	emitter.On("test", func(data interface{}) {
		received = true
	})

	emitter.Emit("test", nil)

	if !received {
		t.Error("expected handler to be called")
	}
}

func TestBasePlugin(t *testing.T) {
	plugin := &BasePlugin{
		NameVal:        "base",
		VersionVal:     "1.0.0",
		DescriptionVal: "A base plugin",
	}

	if plugin.Name() != "base" {
		t.Errorf("expected name 'base', got %s", plugin.Name())
	}
	if plugin.Version() != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", plugin.Version())
	}
	if plugin.Description() != "A base plugin" {
		t.Errorf("expected description 'A base plugin', got %s", plugin.Description())
	}

	err := plugin.Initialize(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = plugin.Shutdown()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

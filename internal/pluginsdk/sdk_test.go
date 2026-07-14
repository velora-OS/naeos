package pluginsdk

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/pluginhost"
)

func TestInterfaceAliases(t *testing.T) {
	var _ pluginhost.Plugin = Plugin(nil)
	var _ pluginhost.Logger = Logger(nil)
	var _ pluginhost.MetricsCollector = MetricsCollector(nil)
	var _ pluginhost.EventEmitter = EventEmitter(nil)
}

func TestStructAliases(t *testing.T) {
	_ = PluginContext{}
	_ = Manifest{}
	_ = ActionManifest{}
	_ = ConfigField{}
	_ = BasePlugin{}
}

func TestPluginStateValues(t *testing.T) {
	if StateCreated != pluginhost.StateCreated {
		t.Errorf("StateCreated mismatch: got %v, want %v", StateCreated, pluginhost.StateCreated)
	}
	if StateInitialized != pluginhost.StateInitialized {
		t.Errorf("StateInitialized mismatch: got %v, want %v", StateInitialized, pluginhost.StateInitialized)
	}
	if StateRunning != pluginhost.StateRunning {
		t.Errorf("StateRunning mismatch: got %v, want %v", StateRunning, pluginhost.StateRunning)
	}
	if StateStopped != pluginhost.StateStopped {
		t.Errorf("StateStopped mismatch: got %v, want %v", StateStopped, pluginhost.StateStopped)
	}
	if StateError != pluginhost.StateError {
		t.Errorf("StateError mismatch: got %v, want %v", StateError, pluginhost.StateError)
	}
}

func TestNewManagerFunc(t *testing.T) {
	mgr := NewManager(t.TempDir())
	if mgr == nil {
		t.Fatal("expected non-nil manager from NewManager()")
	}
}

func TestNewSimpleLoggerFunc(t *testing.T) {
	logger := NewSimpleLogger("test-component")
	if logger == nil {
		t.Fatal("expected non-nil logger from NewSimpleLogger()")
	}
}

func TestNewSimpleMetricsFunc(t *testing.T) {
	metrics := NewSimpleMetrics()
	if metrics == nil {
		t.Fatal("expected non-nil metrics from NewSimpleMetrics()")
	}
}

func TestNewSimpleEventEmitterFunc(t *testing.T) {
	emitter := NewSimpleEventEmitter()
	if emitter == nil {
		t.Fatal("expected non-nil emitter from NewSimpleEventEmitter()")
	}
}

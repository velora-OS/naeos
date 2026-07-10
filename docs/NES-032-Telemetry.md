# NES-032 — Telemetry & Observability

> Status: Draft
> Last Updated: 2026-07-10

Reference for NAEOS telemetry, metrics, and observability.

---

## Overview

NAEOS provides telemetry through the kernel system. Telemetry events track pipeline execution, kernel lifecycle, and component activity.

---

## TelemetryEvent

| Field | Type | Description |
|---|---|---|
| `Name` | `string` | Event name |
| `Timestamp` | `int64` | Unix milliseconds |
| `Payload` | `map[string]any` | Event data |

---

## Metrics

| Field | Type | Description |
|---|---|---|
| `Events` | `int` | Total events emitted |
| `LastEvent` | `TelemetryEvent` | Most recent event |

---

## Kernel Telemetry Events

| Event Name | When | Payload |
|---|---|---|
| `kernel.start` | Before pipeline execution | `services: []string` |
| `kernel.stop` | After pipeline execution | `services: []string` |

---

## Pipeline Telemetry Events

| Event Name | When | Payload |
|---|---|---|
| `pipeline.validate` | After validation | `source_len: int` |
| `pipeline.run` | After full run | `artifacts: int`, `tasks: int`, `reviews: int`, `graph_nodes: int`, `graph_edges: int` |

---

## Usage

### Emit Event

```go
k := kernel.NewKernel()
k.EmitTelemetry(kernel.TelemetryEvent{
    Name:      "custom.event",
    Timestamp: time.Now().UnixMilli(),
    Payload:   map[string]any{"key": "value"},
})
```

### Read Metrics

```go
m := k.Metrics()
fmt.Printf("Total events: %d\n", m.Events)
fmt.Printf("Last event: %s\n", m.LastEvent.Name)
```

### Via CLI

```bash
naeos kernel metrics --config config.yaml
# Output: events=<N> last_event=<name>
```

---

## Runtime Telemetry

The `internal/runtime/telemetry` package provides standalone telemetry:

| Method | Description |
|---|---|
| `Emit(event)` | Emit with auto-timestamp |
| `Events()` | All events |
| `EventsByName(name)` | Filter by name |
| `Metrics()` | Current metrics |
| `Reset()` | Clear all data |

### Metrics Output

```json
{
  "events": 5,
  "start_time": "2026-07-10T10:00:00Z",
  "last_event": {
    "name": "build.complete",
    "timestamp": 1720612800000,
    "payload": {"module": "auth"}
  },
  "event_types": {
    "build.start": 2,
    "build.complete": 2,
    "test.run": 1
  }
}
```

---

## Event Bus Integration

Telemetry events can be published to the kernel event bus:

```go
// Publish telemetry event
k.Publish("telemetry.build", kernel.TelemetryEvent{
    Name:      "build.complete",
    Timestamp: time.Now().UnixMilli(),
    Payload:   map[string]any{"duration_ms": 1234},
})

// Subscribe to telemetry
k.Subscribe("telemetry.build", func(payload any) {
    event := payload.(kernel.TelemetryEvent)
    log.Printf("Event: %s at %d", event.Name, event.Timestamp)
})
```

---

## Observing Pipeline Execution

```bash
# Run with verbose output
naeos --verbose run --config config.yaml --input spec.yaml

# Check metrics after run
naeos kernel metrics --config config.yaml

# List event topics
naeos kernel events --config config.yaml
```

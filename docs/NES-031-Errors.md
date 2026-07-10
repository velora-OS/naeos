# NES-031 — Error Handling Reference

> Status: Draft
> Last Updated: 2026-07-10

Complete reference for NAEOS error codes and troubleshooting.

---

## Error Format

NAEOS errors follow the pattern: `"context: <inner error>"`

```bash
validation failed: project name must not be empty; must contain at least one module
```

---

## Configuration Errors

| Error | Cause | Fix |
|---|---|---|
| `read config: <os error>` | Config file not found or unreadable | Check file path and permissions |
| `config is empty` | Config file is 0 bytes | Write valid config to file |
| `parse config: <error>` | Invalid JSON and YAML | Validate config syntax |
| `missing required --config` | No `--config` flag provided | Add `--config config.yaml` |

---

## Input Errors

| Error | Cause | Fix |
|---|---|---|
| `input cannot be empty` | No input provided | Use `--input` or `--input-file` |
| Both `--input` and `--input-file` | Both flags set | Use only one |

---

## Pipeline Errors

| Error | Cause | Fix |
|---|---|---|
| `input cannot be empty` | Empty spec string | Provide non-empty specification |
| `policy evaluation failed` | Policy rule blocked execution | Check policy rules and context |
| `create artifact dir: <error>` | Cannot create output directory | Check permissions on parent dir |
| `write artifact <path>: <error>` | Cannot write artifact file | Check write permissions |

---

## Kernel Errors

| Error | Cause | Fix |
|---|---|---|
| `service name cannot be empty` | Empty service name | Provide non-empty name |
| `service cannot be nil` | Nil service pointer | Pass non-nil service |
| `service "X" already registered` | Duplicate registration | Use unique names |
| `service "X" not found` | Resolve for unknown service | Register first |
| `kernel already started` | Double Start() | Check lifecycle |
| `initialize service: <error>` | Service Init() failed | Check service implementation |
| `start service: <error>` | Service Start() failed | Check service implementation |
| `kernel is not running` | Stop() before Start() | Start first |
| `stop service: <error>` | Service Stop() failed | Check service implementation |
| `topic cannot be empty` | Empty event topic | Provide topic string |
| `handler cannot be nil` | Nil handler function | Pass valid handler |
| `kernel not initialized` | Kernel is nil | Initialize pipeline first |

---

## Parser Errors

| Error | Cause | Fix |
|---|---|---|
| `input cannot be empty` | Empty spec string | Provide spec content |
| `parse spec: <error>` | Invalid YAML/JSON | Fix syntax |
| `empty specification document` | Parsed YAML is empty | Add content to spec |
| `empty document` | YAML document node empty | Add key-value pairs |
| `map keys must be scalar` | Non-scalar YAML keys | Use string keys |
| `invalid alias node` | Broken YAML alias | Fix alias reference |
| `unsupported YAML node kind` | Unknown node type | Use standard YAML |

---

## Normalizer Errors

| Error | Cause | Fix |
|---|---|---|
| `document is nil` | Nil SpecDocument | Parse first |

---

## Resolver Errors

| Error | Cause | Fix |
|---|---|---|
| `spec is nil` | Nil NormalizedSpec | Normalize first |

---

## Builder Errors

| Error | Cause | Fix |
|---|---|---|
| `resolved spec is nil` | Nil ResolvedSpec | Resolve first |

---

## Validator Errors

| Error | Cause | Fix |
|---|---|---|
| `validation failed: project name must not be empty` | No project name | Set `project` in spec |
| `validation failed: must contain at least one module` | No modules | Add at least one module |
| `validation failed: module[i] name must not be empty` | Module missing name | Set `name` on module |
| `validation failed: module[i] path must not be empty` | Module missing path | Set `path` on module |
| `validation failed: module name "X" is duplicated` | Duplicate module names | Rename duplicate modules |
| `validation failed: service[i] port must be between 0 and 65535` | Invalid port | Use valid port number |
| `metadata.neir_version is recommended` | Missing NEIR version | Add to metadata |

---

## Generator Errors

| Error | Cause | Fix |
|---|---|---|
| `neir is nil` | Nil NEIR model | Build NEIR first |

---

## Runtime Engine Errors

| Error | Cause | Fix |
|---|---|---|
| `artifact is nil` | Nil artifact | Pass valid artifact |
| `artifact path must not be empty` | Empty path string | Set artifact path |
| `no artifacts to execute` | Empty artifact list | Generate artifacts first |
| `failed to execute <path>: <error>` | Execution error | Check artifact content |
| `go file has no content` | Empty .go file | Add Go code |
| `go file missing package declaration` | No `package` line | Add `package X` |
| `yaml file has no content` | Empty .yaml file | Add YAML content |
| `markdown file has no content` | Empty .md file | Add markdown content |

---

## Graph Errors

| Error | Cause | Fix |
|---|---|---|
| `node ID must not be empty` | Empty node ID | Provide ID |
| `node "X" already exists` | Duplicate node | Use unique IDs |
| `node "X" not found` | Unknown node | Add node first |
| `source node "X" not found` | Edge from unknown node | Add source node |
| `target node "X" not found` | Edge to unknown node | Add target node |
| `edge from X to Y already exists` | Duplicate edge | Use unique edges |
| `cycle detected in graph` | Circular dependency | Break the cycle |

---

## Registry Errors

| Error | Cause | Fix |
|---|---|---|
| `entry name must not be empty` | Empty entry name | Provide name |
| `entry "X" already registered` | Duplicate entry | Use unique names |
| `entry "X" not found` | Unknown entry | Register first |

---

## Provenance Errors

| Error | Cause | Fix |
|---|---|---|
| `record ID must not be empty` | Empty record ID | Provide ID |
| `record "X" already exists` | Duplicate record | Use unique IDs |
| `record "X" not found in lineage` | Broken lineage chain | Ensure parent exists |

---

## Policy Errors

| Error | Cause | Fix |
|---|---|---|
| `context is nil` | Nil evaluation context | Pass non-nil context |

---

## Review Errors

| Error | Cause | Fix |
|---|---|---|
| `review input is nil` | Nil review input | Pass valid input |
| `artifact name must not be empty` | Empty artifact name | Provide name |

---

## Troubleshooting Checklist

1. **Config not found** — Check file path exists and is readable
2. **Parse error** — Validate YAML/JSON syntax
3. **Empty output** — Ensure spec has modules and services
4. **Validation fails** — Check all required fields
5. **Policy blocks** — Review policy rules and context values
6. **Artifacts not written** — Check output_dir permissions
7. **Kernel not started** — Ensure Start() called before operations
8. **Cycle detected** — Review module/service dependencies

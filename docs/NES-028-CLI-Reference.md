# NES-028 — CLI Reference

> Status: Draft
> Last Updated: 2026-07-10

Complete reference for the `naeos` CLI tool.

---

## Global Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--verbose` | `bool` | `false` | Enable verbose logging |

```bash
naeos --verbose run --config config.yaml --input spec.yaml
```

---

## Commands

### `naeos version`

Show NAEOS version.

```bash
naeos version
# Output: naeos <version>
```

---

### `naeos init`

Generate a default NAEOS config file.

```bash
naeos init                          # Creates config.example.yaml
naeos init -o myconfig.yaml         # Custom output path
```

| Flag | Short | Type | Default | Description |
|---|---|---|---|---|
| `--output` | `-o` | `string` | `config.example.yaml` | Output path |

**Generated config:**

```yaml
pipeline:
  name: naeos-dev
  mode: development
  verbose: true
  output_dir: ./out
```

---

### `naeos run`

Execute the full NAEOS pipeline.

```bash
naeos run --config config.yaml --input spec.yaml
naeos run --config config.yaml --input "my specification text"
naeos run --config config.yaml --input spec.yaml --output json
naeos run --config config.yaml --input spec.yaml --output yaml --output-file result.yaml
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to JSON or YAML config file |
| `--input` | `string` | `""` | One of `--input`/`--input-file` | Inline specification text |
| `--input-file` | `string` | `""` | One of `--input`/`--input-file` | Path to specification file |
| `--output` | `string` | `text` | No | Output format: `text`, `json`, `yaml` |
| `--output-file` | `string` | `""` | No | Write output to file |
| `--language` | `string` | `""` | No | Target language (repeatable: `--language go --language typescript`) |

**Text output example:**

```
pipeline=my-app mode=development verbose=true output_dir=./out
artifacts=12 tasks=8
```

**JSON output keys:** `pipeline`, `mode`, `verbose`, `output_dir`, `artifacts`, `tasks`

**Multi-language example:**

```bash
# Generate Go + TypeScript artifacts
naeos run --config config.yaml --input spec.yaml --language go --language typescript

# Generate Python artifacts only
naeos run --config config.yaml --input spec.yaml --language python

# All 5 supported languages
naeos run --config config.yaml --input spec.yaml \
  --language go --language typescript --language python \
  --language java --language rust
```

---

### `naeos validate` (alias: `v`)

Validate a specification using the NAEOS pipeline.

```bash
naeos validate --config config.yaml --input spec.yaml
naeos v --config config.yaml --input spec.yaml          # alias
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |
| `--input` | `string` | `""` | One of `--input`/`--input-file` | Inline specification text |
| `--input-file` | `string` | `""` | One of `--input`/`--input-file` | Path to specification file |
| `--output` | `string` | `text` | No | Output format: `text`, `json`, `yaml` |
| `--output-file` | `string` | `""` | No | Write output to file |
| `--language` | `string` | `""` | No | Target language (repeatable) |

**Output keys:** `pipeline`, `mode`, `verbose`, `output_dir`, `status` (`"valid"`), `project`, `source_len`

---

### `naeos inspect`

Inspect the NAEOS pipeline result.

```bash
naeos inspect --config config.yaml --input spec.yaml
naeos inspect --config config.yaml --input spec.yaml --output json
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |
| `--input` | `string` | `""` | Yes | Specification text or file path |
| `--output` | `string` | `text` | No | Output format: `text`, `json`, `yaml` |
| `--output-file` | `string` | `""` | No | Write output to file |

**Output keys:** `pipeline`, `mode`, `verbose`, `output_dir`, `project`, `input`, `artifacts`, `tasks`, `source_words`

---

### `naeos doctor`

Run diagnostics on the NAEOS configuration.

```bash
naeos doctor --config config.yaml
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |

**Output keys:** `config`, `pipeline_name`, `mode`, `output_dir`, `output_dir_exists`

---

### `naeos repair`

Repair the NAEOS output directory (creates directory and README.md if missing).

```bash
naeos repair --config config.yaml
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |

---

### `naeos scaffold`

Generate a starter project scaffold.

```bash
naeos scaffold --name my-project
naeos scaffold --name my-project --output ./scaffold
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--name` | `string` | `""` | Yes | Project name |
| `--output` | `string` | `""` | No | Output directory (defaults to `--name`) |
| `--language` | `string` | `""` | No | Target language (repeatable) |

**Generated structure:**

```
<scaffold>/
  README.md
  spec.yaml
  Makefile
  .gitignore
  Dockerfile
  .github/workflows/ci.yml
  go.mod
  config.yaml
  config.json
  cmd/app/main.go
  internal/core/
    README.md
    package.go
    config.yaml
    handler.go
    repository.go
    service.go
    domain/model.go
    http/handler.go
    http/router.go
    middleware/logging.go
    config/config.go
    config/load.go
    handler_test.go
```

---

### `naeos export`

Export generated artifacts to a directory.

```bash
naeos export --config config.yaml --input spec.yaml
naeos export --config config.yaml --input spec.yaml --output-dir ./artifacts
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |
| `--input` | `string` | `""` | Yes | Specification text or file path |
| `--output-dir` | `string` | `""` | No | Output directory (falls back to config `output_dir`) |
| `--language` | `string` | `""` | No | Target language (repeatable) |

---

### `naeos preview`

Preview generated artifacts without writing them to disk.

```bash
naeos preview --config config.yaml --input spec.yaml
```

| Flag | Type | Default | Required | Description |
|---|---|---|---|---|
| `--config` | `string` | `""` | Yes | Path to config file |
| `--input` | `string` | `""` | Yes | Specification text or file path |

---

### `naeos kernel`

Inspect the NAEOS kernel and service registry. Contains 5 subcommands.

#### Global Kernel Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--config` | `string` | `""` | Path to config file |
| `--output` | `string` | `text` | Output format: `text`, `json`, `yaml` |
| `--topic` | `string` | `""` | Kernel event topic |
| `--payload` | `string` | `""` | Event payload to publish |

#### `naeos kernel services`

List registered kernel services.

```bash
naeos kernel services --config config.yaml
```

#### `naeos kernel metrics`

Show kernel telemetry metrics.

```bash
naeos kernel metrics --config config.yaml
# Output: events=<N> last_event=<name>
```

#### `naeos kernel events`

List active kernel event topics.

```bash
naeos kernel events --config config.yaml
```

#### `naeos kernel publish`

Publish an event to the kernel event bus.

```bash
naeos kernel publish --config config.yaml --topic build.start --payload '{"module":"auth"}'
```

#### `naeos kernel subscribe`

Subscribe to a kernel event topic.

```bash
naeos kernel subscribe --config config.yaml --topic build.start
naeos kernel subscribe --config config.yaml --topic build.start --payload '{"test":"data"}'
```

---

## Common Workflows

### Quick Start

```bash
# 1. Generate config
naeos init -o config.yaml

# 2. Create spec.yaml (or use scaffold)
naeos scaffold --name my-app

# 3. Validate spec
naeos validate --config config.yaml --input-file my-app/spec.yaml

# 4. Run pipeline
naeos run --config config.yaml --input-file my-app/spec.yaml

# 5. Export artifacts
naeos export --config config.yaml --input-file my-app/spec.yaml
```

### Diagnostics

```bash
naeos doctor --config config.yaml        # Check config health
naeos repair --config config.yaml        # Fix output directory
```

### Kernel Inspection

```bash
naeos kernel services --config config.yaml    # List services
naeos kernel events --config config.yaml      # List topics
naeos kernel metrics --config config.yaml     # View metrics
```

### Multi-Language SDK Generation

```bash
# 1. Validate with language targets
naeos validate --config config.yaml --input-file spec.yaml --language go --language typescript

# 2. Run pipeline with multi-language output
naeos run --config config.yaml --input-file spec.yaml --language go --language typescript --language python

# 3. Export artifacts
naeos export --config config.yaml --input-file spec.yaml --language go --language typescript --language python

# 4. Preview what would be generated
naeos preview --config config.yaml --input-file spec.yaml --language go
```

---

## Output Formats

All commands that accept `--output` support:

| Format | Description |
|---|---|
| `text` | Human-readable key=value pairs (default) |
| `json` | Pretty-printed JSON |
| `yaml` | YAML output |

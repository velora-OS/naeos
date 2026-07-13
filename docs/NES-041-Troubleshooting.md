# NES-041 — Troubleshooting Guide

> Status: Draft
> Last Updated: 2026-07-13

Comprehensive troubleshooting scenarios for NAEOS operations.

> For the full error code reference, see [NES-031-Errors.md](NES-031-Errors.md).

---

## 1. Pipeline Issues

### 1.1 Pipeline hangs on parse

**Symptoms:** `naeos run` never completes, process stuck at parser stage.

**Causes:**
- Circular `$ref` references in the specification
- Extremely large spec files (>10MB)

**Fix:**

```bash
# Enable verbose output to identify the stuck stage
naeos --verbose run --config config.yaml --input spec.yaml

# Validate spec structure first
naeos validate --config config.yaml --input-file spec.yaml

# Check for circular references
naeos --verbose run --config config.yaml --input spec.yaml 2>&1 | grep -i cycle
```

If the spec is large, split it into multiple modules:

```yaml
# Instead of one large module, use separate modules
modules:
  - name: auth
    path: ./internal/auth
  - name: user
    path: ./internal/user
```

### 1.2 Validation fails with "dangling dependency"

**Symptoms:** `validation failed: module dependency "X" not found`

**Cause:** A module references another module that doesn't exist in the spec.

**Fix:**

```bash
# Run validation to see the full error
naeos validate --config config.yaml --input-file spec.yaml

# Check all module names match
naeos validate --config config.yaml --input-file spec.yaml --output json
```

Ensure all referenced modules are defined:

```yaml
modules:
  - name: auth
    path: ./internal/auth
    dependencies:
      - user        # Must exist as a module above
      - config      # Must exist as a module above
  - name: user
    path: ./internal/user
  - name: config
    path: ./internal/config
```

### 1.3 Export generates empty output

**Symptoms:** `naeos export` produces no artifacts.

**Causes:**
- Missing `project` name in spec
- No modules or services defined

**Fix:**

```bash
# Check what the pipeline produces
naeos inspect --config config.yaml --input-file spec.yaml

# Ensure project has content
naeos validate --config config.yaml --input-file spec.yaml
```

Minimum valid spec:

```yaml
project: my-project
modules:
  - name: core
    path: ./internal/core
```

### 1.4 Graph cycle detected

**Symptoms:** `cycle detected in graph`

**Cause:** Circular dependency between modules or services.

**Fix:**

```bash
# Verbose output shows the cycle path
naeos --verbose run --config config.yaml --input-file spec.yaml 2>&1 | grep cycle
```

Restructure dependencies to break the cycle:

```yaml
# Before (cycle: A -> B -> A)
modules:
  - name: A
    dependencies: [B]
  - name: B
    dependencies: [A]

# After (introduce interface module)
modules:
  - name: A
    dependencies: [iface]
  - name: B
    dependencies: [iface]
  - name: iface
    path: ./internal/iface
```

---

## 2. Cloud Deployment Issues

### 2.1 Deploy fails with region error

**Symptoms:** `deploy failed: invalid region "X"`

**Cause:** Region format doesn't match the cloud provider.

**Fix:**

```bash
# Check provider-specific region format
naeos cloud plan --config config.yaml --input-file spec.yaml --provider aws --region us-east-1
naeos cloud plan --config config.yaml --input-file spec.yaml --provider gcp --region us-central1
naeos cloud plan --config config.yaml --input-file spec.yaml --provider azure --region eastus
```

Valid region formats:

| Provider | Format | Examples |
|----------|--------|---------|
| AWS | `us-<direction>-<number>` | `us-east-1`, `eu-west-2` |
| GCP | `<region>-<zone>` | `us-central1`, `asia-east1-a` |
| Azure | `<region><letter>` | `eastus`, `westeurope` |
| DigitalOcean | `<city>-<number>` | `nyc1`, `sfo3` |

### 2.2 Terraform init fails

**Symptoms:** `terraform init` error during cloud operations.

**Causes:**
- Terraform binary not installed
- Terraform not in PATH
- Version too old

**Fix:**

```bash
# Check terraform installation
terraform version

# If not installed, install >= 1.5
# macOS
brew install terraform

# Linux
wget https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip
unzip terraform_1.5.7_linux_amd64.zip
sudo mv terraform /usr/local/bin/

# Verify
terraform version
```

### 2.3 State file corrupted

**Symptoms:** `failed to read state: invalid JSON` or concurrent deployment errors.

**Fix:**

```bash
# Find the state file
ls ~/.naeos/cloud/<project>/state.json

# Backup and remove corrupted state
cp ~/.naeos/cloud/<project>/state.json ~/.naeos/cloud/<project>/state.json.bak
rm ~/.naeos/cloud/<project>/state.json

# Redeploy
naeos deploy --config config.yaml --input-file spec.yaml --environment staging
```

### 2.4 Deploy dry-run shows unexpected resources

**Symptoms:** Resources differ from what was planned.

**Fix:**

```bash
# Generate fresh plan
naeos cloud plan --config config.yaml --input-file spec.yaml --provider aws --region us-east-1 --output json

# Compare with dry-run
naeos deploy --config config.yaml --input-file spec.yaml --environment staging --dry-run --output json
```

---

## 3. Plugin Issues

### 3.1 Plugin load fails on Alpine

**Symptoms:** `plugin load failed: not a valid Win32 application` or `exec format error`

**Cause:** Plugin was compiled with glibc dependencies, incompatible with Alpine's musl.

**Fix:**

```bash
# Rebuild with static linking
CGO_ENABLED=0 go build -o my-plugin.so -buildmode=c-shared ./cmd/plugin

# Or use Alpine-compatible base in Dockerfile
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev
```

### 3.2 Plugin timeout

**Symptoms:** `plugin execution timed out after 30s`

**Fix:**

Increase the timeout in the plugin manifest:

```yaml
# plugin.yaml
name: my-plugin
timeout: 60s    # default is 30s
```

Or optimize the plugin code:

```bash
# Profile the plugin
naeos plugin test --source ./my-plugin --output json

# Check duration
# duration_ms: 45000 (exceeds default 30s)
```

### 3.3 WASM plugin instantiation fails

**Symptoms:** `wasm instantiation failed: invalid module`

**Cause:** WASM binary is not WASI-compliant.

**Fix:**

```bash
# Validate the WASM binary
wasm-tools validate my-plugin.wasm --features=component-model

# Rebuild with WASI target
rustup target add wasm32-wasi
cargo build --target wasm32-wasi --release

# Or with Go
GOOS=wasip1 GOARCH=wasm go build -o my-plugin.wasm ./cmd/plugin
```

### 3.4 Plugin not found after install

**Symptoms:** `plugin "X" not found in registry`

**Fix:**

```bash
# Test plugin first
naeos plugin test --source ./my-plugin

# Check plugin directory
ls ~/.naeos/plugins/

# Reinstall if needed
curl -X POST http://localhost:8080/api/v1/plugins \
  -d '{"name":"my-plugin","source":"./my-plugin"}'
```

---

## 4. API Server Issues

### 4.1 Pipeline run returns 429

**Symptoms:** `HTTP 429: rate limit exceeded`

**Cause:** Too many requests without authentication.

**Fix:**

```bash
# Use API key for higher rate limits
curl -H "Authorization: Bearer $NAEOS_API_KEY" \
  -X POST http://localhost:8080/api/v1/pipeline/run \
  -d '{"config":"config.yaml","input":"spec.yaml"}'

# Or wait for rate limit to reset
sleep 60
```

### 4.2 WebSocket connection drops

**Symptoms:** `ws connection closed: 1006 abnormal closure`

**Fix:**

```bash
# Check proxy settings
echo $HTTP_PROXY
echo $HTTPS_PROXY

# Increase timeout in config
naeos config set --config config.yaml --key api.ws_timeout --value 300s

# Test connection
naeos events stream --config config.yaml --topic build.*
```

### 4.3 JWT auth fails

**Symptoms:** `unauthorized: token expired or invalid`

**Fix:**

```bash
# Regenerate token
curl -X POST http://localhost:8080/api/v1/auth/token \
  -d '{"user":"admin","password":"..."}'

# Set in environment
export NAEOS_API_KEY="eyJhbG..."

# Verify token
naeos health --config config.yaml
```

---

## 5. General Issues

### 5.1 Build fails with "go.sum mismatch"

**Symptoms:** `go.sum mismatch` or `go: verifying module checksums`

**Fix:**

```bash
# Clean and rebuild dependencies
go mod tidy
go mod download

# If cache is corrupted
go clean -modcache
go mod tidy
go mod download

# Full rebuild
make clean && make build
```

### 5.2 Version shows 0.0.0

**Symptoms:** `naeos version` outputs `naeos 0.0.0`

**Cause:** VERSION file missing or build didn't inject version.

**Fix:**

```bash
# Check VERSION file exists
cat VERSION

# If missing, create it
echo "0.1.0" > VERSION

# Rebuild with version injection
make build

# Verify
naeos version
# Output: naeos 0.1.0
```

### 5.3 Config reload not working

**Symptoms:** Changes to config file not reflected.

**Fix:**

```bash
# Check file permissions
ls -la config.yaml

# Validate config syntax
naeos config validate --config config.yaml

# Show resolved config to verify changes
naeos config show --config config.yaml

# Run diagnostics
naeos doctor --config config.yaml
```

### 5.4 Output directory not writable

**Symptoms:** `create artifact dir: permission denied`

**Fix:**

```bash
# Check output directory permissions
ls -la $(dirname config.yaml)

# Repair output directory
naeos repair --config config.yaml

# Create manually if needed
mkdir -p ./out
chmod 755 ./out
```

### 5.5 Specification parse error

**Symptoms:** `parse spec: yaml: line X: did not find expected...`

**Fix:**

```bash
# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('spec.yaml'))"

# Check for common issues:
# - Missing quotes around strings with special characters
# - Incorrect indentation (use spaces, not tabs)
# - Trailing colons or commas

# Minimal valid spec
cat << 'EOF' > spec.yaml
project: my-project
modules:
  - name: core
    path: ./internal/core
EOF

naeos validate --config config.yaml --input-file spec.yaml
```

### 5.6 Kernel not running

**Symptoms:** `kernel is not running`

**Fix:**

```bash
# Ensure pipeline is started before kernel operations
naeos run --config config.yaml --input-file spec.yaml

# Then use kernel commands
naeos kernel services --config config.yaml
naeos kernel metrics --config config.yaml
```

---

## 6. Diagnostic Commands

Use these commands to gather information before filing issues:

```bash
# Full system check
naeos doctor --config config.yaml

# Health check with details
naeos health --config config.yaml --output json

# Pipeline inspection
naeos inspect --config config.yaml --input-file spec.yaml --output json

# Config validation
naeos config validate --config config.yaml

# Version check
naeos version

# Kernel state
naeos kernel services --config config.yaml
naeos kernel metrics --config config.yaml
naeos kernel events --config config.yaml
```

---

## 7. Getting Help

- Check [NES-031-Errors.md](NES-031-Errors.md) for the complete error code reference
- Check [NES-028-CLI-Reference.md](NES-028-CLI-Reference.md) for command documentation
- Run `naeos --help` for available commands
- Run `naeos <command> --help` for command-specific options

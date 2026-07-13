# NES-030 — Specification Language

> Status: Stable
> Last Updated: 2026-07-13

Complete reference for the NAEOS Specification Language (spec.yaml).

---

## Overview

The NAEOS specification is a YAML or JSON document that describes a software project's structure, services, architecture, deployment, and testing strategy. It serves as the single source of truth for code generation.

---

## Root Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `project` | `string` | No | Project name (auto-generated if omitted) |
| `modules` | `[]Module` | No | Application modules |
| `services` | `[]Service` | No | Service definitions |
| `architecture` | `Architecture` | No | Architecture pattern |
| `deployment` | `Deployment` | No | Deployment strategy |
| `testing` | `Testing` | No | Testing strategy |
| `cloud` | `Cloud` | No | Cloud deployment configuration |
| `plugins` | `[]Plugin` | No | Plugin extensions |
| `ai` | `AI` | No | AI integration settings |

---

## Module

Defines a code module in the project.

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | `string` | Yes | Module name |
| `path` | `string` | Yes | File system path |
| `description` | `string` | No | Module description |
| `dependencies` | `[]string` | No | Other modules this depends on |

```yaml
modules:
  - name: auth
    path: ./internal/auth
    description: Authentication and authorization module
    dependencies:
      - core
      - user
```

---

## Service

Defines a runnable service.

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | `string` | Yes | Service name |
| `kind` | `string` | Yes | Service type: `http`, `grpc`, `worker`, `cli`, `job` |
| `port` | `int` | Yes | Listening port |
| `description` | `string` | No | Service description |
| `endpoints` | `[]Endpoint` | No | API endpoints |

### Endpoint

| Field | Type | Required | Description |
|---|---|---|---|
| `method` | `string` | Yes | HTTP method: `GET`, `POST`, `PUT`, `DELETE`, `PATCH` |
| `path` | `string` | Yes | URL path |
| `action` | `string` | Yes | Handler action name |

```yaml
services:
  - name: api
    kind: http
    port: 8080
    description: Main API server
    endpoints:
      - method: GET
        path: /health
        action: healthCheck
      - method: POST
        path: /users
        action: createUser
      - method: GET
        path: /users/:id
        action: getUser
```

---

## Architecture

Defines the architectural pattern.

| Field | Type | Required | Description |
|---|---|---|---|
| `pattern` | `string` | Yes | Architecture pattern (see values below) |
| `description` | `string` | No | Pattern description |
| `principles` | `[]string` | No | Design principles |

### Supported Patterns

| Pattern | Description |
|---|---|
| `layered` | Traditional N-tier layered architecture |
| `clean` | Clean Architecture (Uncle Bob) |
| `hexagonal` | Hexagonal / Ports & Adapters |
| `microkernel` | Plugin-based microkernel |
| `event-driven` | Event-driven architecture |
| `cqrs` | Command Query Responsibility Segregation |
| `monolith` | Single deployable unit |

```yaml
architecture:
  pattern: hexagonal
  description: Ports and adapters architecture
  principles:
    - Dependency inversion
    - Separation of concerns
    - Testability
```

---

## Deployment

Defines deployment strategy and environments.

| Field | Type | Required | Description |
|---|---|---|---|
| `strategy` | `string` | Yes | Deployment strategy (see values below) |
| `environments` | `[]string` | Yes | Target environments |

### Supported Strategies

| Strategy | Description |
|---|---|
| `rolling` | Rolling update deployment |
| `blue-green` | Blue-green deployment |
| `canary` | Canary release |
| `recreate` | Full stop then recreate |

```yaml
deployment:
  strategy: blue-green
  environments:
    - staging
    - production
```

---

## Testing

Defines testing strategy and coverage targets.

| Field | Type | Required | Description |
|---|---|---|---|
| `strategy` | `string` | Yes | Testing strategy (see values below) |
| `coverage` | `string` | No | Minimum coverage target (e.g., `"80%"`) |

### Supported Strategies

| Strategy | Description |
|---|---|
| `unit` | Unit testing focus |
| `integration` | Integration testing focus |
| `e2e` | End-to-end testing |
| `contract` | Contract testing |

```yaml
testing:
  strategy: unit
  coverage: "80%"
```

---

## Cloud

Defines cloud provider deployment configuration.

| Field | Type | Required | Description |
|---|---|---|---|
| `provider` | `string` | Yes | Cloud provider: `aws`, `gcp`, `azure`, `digitalocean` |
| `region` | `string` | No | Target region |
| `services` | `[]CloudService` | No | Cloud-managed services to provision |
| `scaling` | `Scaling` | No | Auto-scaling configuration |

### CloudService

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | `string` | Yes | Service instance name |
| `type` | `string` | Yes | Service type: `compute`, `storage`, `database`, `cache`, `queue` |
| `tier` | `string` | No | Service tier: `small`, `medium`, `large` |

### Scaling

| Field | Type | Required | Description |
|---|---|---|---|
| `min_replicas` | `int` | No | Minimum replicas |
| `max_replicas` | `int` | No | Maximum replicas |
| `target_cpu` | `string` | No | CPU utilization target (e.g., `"70%"`) |

```yaml
cloud:
  provider: aws
  region: us-east-1
  services:
    - name: api-server
      type: compute
      tier: medium
    - name: main-db
      type: database
      tier: small
    - name: session-cache
      type: cache
      tier: small
  scaling:
    min_replicas: 2
    max_replicas: 10
    target_cpu: "70%"
```

---

## Plugins

Defines plugin extensions for the pipeline.

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | `string` | Yes | Plugin name |
| `source` | `string` | Yes | Plugin source: registry name, local path, or URL |
| `version` | `string` | No | Plugin version constraint |
| `config` | `map` | No | Plugin-specific configuration |

```yaml
plugins:
  - name: terraform-generator
    source: naeos-terraform
    version: ">=1.0.0"
    config:
      provider: aws
      state_backend: s3
  - name: custom-validator
    source: ./plugins/custom-validator
    config:
      rules:
        - require-description
        - max-endpoints-50
```

---

## AI

Defines AI integration settings for context generation and enrichment.

| Field | Type | Required | Description |
|---|---|---|---|
| `context_type` | `string` | No | Context bundle type: `full`, `summary`, `dependencies` |
| `enrichment` | `[]string` | No | Enrichment focus areas: `security`, `performance`, `testing` |
| `mcp_tools` | `[]MCPTool` | No | Custom MCP tool definitions |

### MCPTool

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | `string` | Yes | Tool name |
| `description` | `string` | Yes | Tool description |
| `parameters` | `map` | No | JSON Schema for parameters |

```yaml
ai:
  context_type: full
  enrichment:
    - security
    - performance
  mcp_tools:
    - name: deploy_preview
      description: Deploy a preview environment for the spec
      parameters:
        type: object
        properties:
          environment:
            type: string
```

---

## Full Example

```yaml
project: my-api

modules:
  - name: auth
    path: ./internal/auth
    description: Authentication module
    dependencies:
      - core
  - name: core
    path: ./internal/core
    description: Core business logic
  - name: user
    path: ./internal/user
    description: User management
    dependencies:
      - core

services:
  - name: api
    kind: http
    port: 8080
    description: Main REST API
    endpoints:
      - method: GET
        path: /health
        action: healthCheck
      - method: POST
        path: /auth/login
        action: login
      - method: POST
        path: /auth/register
        action: register
      - method: GET
        path: /users
        action: listUsers
      - method: GET
        path: /users/:id
        action: getUser

architecture:
  pattern: hexagonal
  description: Hexagonal architecture for testability
  principles:
    - Core logic independent of frameworks
    - Ports define interfaces
    - Adapters implement ports

deployment:
  strategy: rolling
  environments:
    - development
    - staging
    - production

testing:
  strategy: unit
  coverage: "85%"

cloud:
  provider: aws
  region: us-east-1
  services:
    - name: api-server
      type: compute
      tier: medium
    - name: main-db
      type: database
      tier: small
  scaling:
    min_replicas: 2
    max_replicas: 8
    target_cpu: "70%"

plugins:
  - name: terraform-generator
    source: naeos-terraform
    version: ">=1.0.0"
    config:
      provider: aws

ai:
  context_type: full
  enrichment:
    - security
    - performance
```

---

## Auto-Defaults

When fields are omitted, the parser applies defaults:

| Condition | Default |
|---|---|
| `project` is empty | Slugified from input or `"default-project"` |
| `modules` is empty | Single module: name=`"default-module"`, path=`"./default"` |
| `module.path` is empty | `"./<slugified-name>"` |
| `cloud.region` is empty | Provider's default region |
| `cloud.scaling` is empty | `min_replicas=1`, `max_replicas=3`, `target_cpu="80%"` |

---

## Format Support

| Format | Parser |
|---|---|
| YAML | `gopkg.in/yaml.v3` |
| JSON | `encoding/json` |

Both formats are accepted. The parser tries JSON first, then YAML.

---

## Validation

The NEIR validator checks:

| Rule | Error |
|---|---|
| Project name required | `"project name must not be empty"` |
| At least one module | `"must contain at least one module"` |
| Module name required | `"module[i] name must not be empty"` |
| Module path required | `"module[i] path must not be empty"` |
| No duplicate module names | `"module name is duplicated"` |
| Service port 0-65535 | `"port must be between 0 and 65535"` |
| Cloud provider required | `"cloud.provider must be specified when cloud section is present"` |
| Plugin source required | `"plugin[i] source must not be empty"` |

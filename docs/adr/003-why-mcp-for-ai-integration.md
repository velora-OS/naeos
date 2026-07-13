# ADR-003: Why MCP for AI Integration

> Status: Accepted
> Date: 2026-07-10

## Context

The AI coding assistant ecosystem is growing rapidly with tools like Claude, Copilot, Cursor, and others. Each tool has its own integration protocol. NAEOS needs to expose its capabilities (spec reading, code generation, validation, context bundling) to AI agents without coupling to any single vendor's API.

## Decision

Use the **Model Context Protocol (MCP)** as the standardized interface for AI agent integration.

## Consequences

### Positive

- Single MCP server exposes NAEOS tools to any MCP-compatible AI agent
- Protocol is transport-agnostic (stdio, HTTP, SSE) and works across environments
- Tool definitions are self-describing, enabling AI agents to discover NAEOS capabilities at runtime
- Reduces integration surface: one server instead of per-vendor adapters

### Negative

- MCP is still evolving; protocol changes may require updates to the server
- Some AI tools do not yet support MCP natively, requiring fallback adapters
- Debugging MCP interactions requires JSON-RPC level inspection

### Mitigations

- The MCP server is isolated behind `/mcp/message` in the API and as a standalone process
- Fallback REST endpoints remain available for non-MCP clients
- MCP server version is pinned per release and updated deliberately
- `naeos doctor` validates MCP server configuration health

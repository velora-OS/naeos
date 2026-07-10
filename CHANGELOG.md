# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Documentation index with recommended reading orders (beginner, policy, profile, CLI, testing).
- NES-028 CLI Reference — comprehensive CLI command documentation.
- NES-029 Configuration — pipeline configuration reference.
- NES-030 Specification Language — NAEOS specification language docs.
- NES-031 Errors — exhaustive error catalog.
- NES-032 Telemetry — telemetry and metrics reference.
- NES-033 Testing Guide — test guide with coverage requirements.
- NES-034 Event Bus — internal pub/sub event bus documentation.
- NES-035 Version Management — SemVer management documentation.
- NES-036 Template Renderer — template rendering engine documentation.
- NES-037 Knowledge Graph & Provenance — knowledge graph and lineage documentation.
- NES-038 Shared Types & Contracts — shared types and contracts documentation.
- NAEOS-GOV-002 Vision — long-term vision document.
- NAEOS-GOV-005 Core Principles — 8 core engineering principles.
- Expanded 18 NES stub documents (NES-003 through NES-022) with full API references and examples.

### Changed
- NES-001 Repository — updated repository structure to match actual codebase paths.
- DOCUMENTATION-INDEX.md — added NES-028 through NES-038, Go package reference section, CLI and testing reading orders.

## [0.1.0] - 2026-01-01

### Added
- Initial project structure.
- CLI with 11 subcommands: init, run, validate, inspect, doctor, repair, scaffold, export, preview, kernel, version.
- Core pipeline: parser, normalizer, resolver, NEIR builder, validator.
- Planner: DAG graph with topological sort and cycle detection.
- Generator engine: Go project code, Dockerfile, CI, documentation.
- Policy evaluator with 7 operators and 5 default rules.
- Artifact reviewer with governance rules.
- Knowledge graph with 14 node types and 13 edge types.
- Provenance tracking store.
- Runtime execution engine with deduplication.
- Telemetry event collector.
- 34 modular design documents (NES-000 through NES-033).
- 10 specification documents (NAEOS-SPEC-001 through 010).
- 8 constitutional documents (NAEOS-CON-001 through 008).
- 8 governance documents (NAEOS-GOV-001 through 008).
- 4 kernel specification documents (NAEOS-KER-001 through 004).
- 7 policy system documents (NAEOS-POL-001 through 007).
- 7 profile system documents (NAEOS-PRO-001 through 007).
- 1 reference architecture document (NAEOS-NRA-001).
- ADR and RFC templates with examples.
- Example specifications (minimal and full).

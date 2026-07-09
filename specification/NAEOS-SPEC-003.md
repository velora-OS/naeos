Document ID: NAEOS-SPEC-003

Title: Universal Artifact Model

Short Name: UAM

Version: 1.0.0

Status: Stable

Category: Core Specification

Normative: true

Priority: Critical

Owner: NAEOS Foundation

Depends On:

- SPEC-001

- SPEC-002
Universal Artifact Model (UAM)
Executive Summary

Universal Artifact Model (UAM) adalah model data universal yang digunakan oleh seluruh artefak dalam ekosistem NAEOS.

Setiap artefak direpresentasikan sebagai Artifact Object yang memiliki struktur metadata, konten, relasi, versi, dan jejak perubahan yang konsisten.

Dengan UAM, compiler hanya perlu memahami satu model data untuk memproses seluruh jenis artefak.

1. Purpose

UAM bertujuan untuk:

Menyatukan representasi seluruh artefak.
Menyederhanakan compiler dan validator.
Memastikan interoperabilitas.
Mendukung traceability dan versioning.
Mempermudah integrasi dengan AI.
2. Artifact Definition

Sebuah Artifact adalah setiap objek engineering yang memiliki identitas, metadata, konten, dan hubungan dengan artefak lain.

Contoh Artifact:

Document
Standard
Constitution
Rule
RFC
ADR
Template
Prompt
Workflow
Source Code
Test Case
API Contract
Deployment Manifest
3. Artifact Lifecycle
4. Universal Structure
artifact:

  id:

  type:

  title:

  summary:

  owner:

  category:

  version:

  status:

  created_at:

  updated_at:

  tags:

  labels:

  references:

  dependencies:

  content:
5. Mandatory Fields

Semua Artifact MUST memiliki:

Field	Requirement
id	MUST
type	MUST
title	MUST
version	MUST
status	MUST
owner	MUST
content	MUST
6. Artifact Types

Jenis artefak standar:

Document
Specification
Constitution
Standard
Policy
Rule
RFC
ADR
Template
Playbook
Guide
Checklist
Schema
Prompt
Workflow
API
Component
Module
Package
Repository
Test
Deployment
Knowledge

Implementasi dapat menambahkan tipe baru tanpa mengubah spesifikasi inti.

7. Identity Model

Setiap Artifact memiliki Artifact ID yang unik.

Format:

<DOMAIN>-<CATEGORY>-<NUMBER>

Contoh:

NAEOS-GOV-001
NAEOS-SPEC-003
NAEOS-STD-012
NAEOS-RFC-0042
NAEOS-ADR-0015

Artifact ID MUST NOT berubah selama siklus hidup artefak.

8. Metadata Model

Metadata dibagi menjadi:

Identity Metadata
Ownership Metadata
Version Metadata
Classification Metadata
Relationship Metadata
Security Metadata

Metadata harus dapat diproses mesin tanpa membaca isi dokumen.

9. Relationship Model

Setiap Artifact dapat memiliki relasi:

depends_on

references

extends

implements

supersedes

duplicates

derived_from

generated_by

owned_by

Relasi harus konsisten dengan Engineering Knowledge Graph.

10. Traceability

Setiap Artifact harus dapat ditelusuri.

Contoh:

Business Requirement

↓

Specification

↓

Architecture

↓

Component

↓

Source Code

↓

Test

↓

Deployment

Compiler harus mampu membangun rantai jejak ini secara otomatis.

11. Version Model

Artifact mengikuti Semantic Versioning.

Perubahan metadata tidak selalu mengubah versi.

Perubahan normatif MUST menaikkan versi sesuai kebijakan Versioning Policy.

12. Validation Rules

Validator harus memeriksa:

ID unik.
Metadata lengkap.
Status valid.
Relasi valid.
Referensi tidak rusak.
Dependensi tidak melingkar (kecuali diizinkan).
13. Serialization

Artifact harus dapat diekspor ke:

Markdown
YAML
JSON
XML (opsional)
PDF (render)
HTML
Graph Node

Tanpa kehilangan informasi penting.

14. Compiler Behavior

Compiler memperlakukan semua Artifact dengan alur yang sama:

Artifact

↓

Parser

↓

Validator

↓

Knowledge Graph

↓

Transformation Engine

↓

Output Adapter

Perbedaan hanya terletak pada adapter output.

15. Extension Model

Vendor atau organisasi dapat membuat Artifact baru.

Contoh:

Healthcare Guideline

extends

Standard

Ekstensi tidak boleh mengubah perilaku Artifact inti.

16. Conformance

Implementasi UAM MUST:

menggunakan struktur Artifact resmi,
mendukung metadata wajib,
mendukung traceability,
mendukung relationship,
kompatibel dengan Engineering Knowledge Graph.
17. Related Documents
ID	Document
NAEOS-SPEC-001	Overview
NAEOS-SPEC-002	Engineering Knowledge Graph
NAEOS-SPEC-004	Metadata Specification
NAEOS-SPEC-005	Rule Model
NAEOS-GOV-008	Versioning Policy
Revision History
Version	Date	Change
1.0.0	2026-07-09	Initial Universal Artifact Model
Status
NAEOS-SPEC-003

APPROVED

Universal Artifact Model Established

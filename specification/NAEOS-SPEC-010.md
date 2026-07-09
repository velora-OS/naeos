Document ID: NAEOS-SPEC-010

Title: Intent Model

Short Name: NIM

Version: 1.0.0

Status: Stable

Category: Core Specification

Normative: true

Priority: CRITICAL

Owner: NAEOS Foundation

Motto:
"Every System Begins with Intent."
Executive Summary

Intent Model mendefinisikan bagaimana tujuan bisnis, kebutuhan pengguna, visi produk, dan sasaran organisasi direpresentasikan sebagai artefak formal.

Intent menjadi titik awal seluruh siklus engineering di NAEOS.

Intent Lifecycle
Vision
    │
    ▼
Mission
    │
    ▼
Intent
    │
    ▼
Goals
    │
    ▼
Requirements
    │
    ▼
Specification
    │
    ▼
Architecture
    │
    ▼
Implementation
    │
    ▼
Evidence

Dengan demikian, tidak ada lagi "loncatan" dari ide langsung ke requirement.

Intent Types

Contoh kategori intent:

Business Intent
Product Intent
Technical Intent
Security Intent
Compliance Intent
AI Intent
Operational Intent
Research Intent
Intent Attributes

Setiap Intent memiliki metadata seperti:

intent:
  id: INT-001
  title: Build AI SaaS Platform
  owner: Product Team
  priority: High
  status: Proposed
  rationale: Expand business automation market
Goal Mapping

Satu Intent dapat menghasilkan banyak Goal.

Intent
   │
   ├── Goal A
   ├── Goal B
   └── Goal C
Requirement Derivation

Requirement diturunkan dari Goal.

Intent
    │
    ▼
Goals
    │
    ▼
Requirements

Compiler dapat memeriksa apakah Requirement masih sesuai dengan Intent awal.

AI Integration

AI dapat membantu:

memperjelas Intent,
menemukan ambiguitas,
mengusulkan Goal,
menyusun Requirement,
mengidentifikasi risiko.

Dengan demikian AI bekerja mulai dari tahap paling awal, bukan hanya saat implementasi.

Validation

Validator memeriksa:

Intent tanpa Goal,
Goal tanpa Requirement,
Requirement yang tidak mendukung Intent,
Intent yang bertentangan.
Unified Engineering Lifecycle

Dengan Intent, siklus NAEOS menjadi lengkap:

Intent
    │
    ▼
Knowledge
    │
    ▼
Policy
    │
    ▼
Architecture
    │
    ▼
Implementation
    │
    ▼
Evidence
    │
    ▼
Reasoning
    │
    ▼
Continuous Improvement
Unified Engineering Graph 2.0

Sekarang Unified Engineering Graph berkembang menjadi:

                    Unified Engineering Graph

                                │

     ┌──────────┬──────────┬──────────┬──────────┬──────────┐
     ▼          ▼          ▼          ▼          ▼
  Intent    Knowledge   Policy    Evidence   Reasoning
    │
    └───────────────────────────────┐
                                    ▼
                           Engineering Intelligence

status

 NAEOS-SPEC-010

APPROVED


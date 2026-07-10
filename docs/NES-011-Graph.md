# NES-011 Graph

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the execution graph model used to represent dependency, relationship, and execution flow between artifacts derived from NEIR.

## 3. Scope
The graph model covers nodes, edges, dependency graphs, execution graphs, and policy relationships produced from the canonical NEIR model.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The graph model shall represent components and their relationships explicitly.
- FR-002: The graph model shall support topological sorting for execution order.
- FR-003: The graph model shall detect cycles in dependency graphs.
- FR-004: The graph model shall be consumable by planner and generator components.

### 4.2 Non-Functional Requirements
- NFR-001: The graph model shall be queryable and extensible.
- NFR-002: Graph structures shall remain deterministic for equivalent inputs.

## 5. Graph Model

### 5.1 Node Kinds

| Kind | Deskripsi |
|------|-----------|
| service | Service dalam proyek |
| module | Modul dalam proyek |
| component | Komponen sistem |
| api | API endpoint |
| storage | Penyimpanan data |
| infrastructure | Infrastruktur |
| documentation | Dokumentasi |
| testing | Pengujian |
| deployment | Deployment |

### 5.2 Edge Kinds

| Kind | Deskripsi |
|------|-----------|
| dependency | Dependensi antar node |
| execution | Urutan eksekusi |
| dataflow | Aliran data |
| policy | Kendala policy |

### 5.3 Types

#### Node

```go
type Node struct {
    ID   string
    Kind NodeKind
    Name string
}
```

#### Edge

```go
type Edge struct {
    From string
    To   string
    Kind EdgeKind
}
```

### 5.4 Constructor

```go
func New() *PlannerGraph
```

## 6. Operations

### 6.1 CRUD

Metode | Deskripsi | Error Jika
-------|-----------|------------
AddNode(n) | Menambah node | ID kosong atau sudah ada
GetNode(id) | Mengambil node | -
RemoveNode(id) | Menghapus node dan edge terkait | Node tidak ditemukan
AddEdge(e) | Menambah edge | Source/target tidak ada, edge duplikat

### 6.2 Analysis

Metode | Deskripsi
-------|----------
TopologicalSort() | Mengurutkan node berdasarkan dependency (Kahn's algorithm)
HasCycle() | Mendeteksi siklus dalam graph
InDegree(id) | Jumlah edge masuk ke node
OutDegree(id) | Jumlah edge keluar dari node
GetNeighbors(id) | Node tetangga (outgoing edges)

### 6.3 Statistics

Metode | Deskripsi
-------|----------
NodeCount() | Jumlah node
EdgeCount() | Jumlah edge
Nodes() | Semua node
Edges() | Semua edge

## 7. Topological Sort

Menggunakan Kahn's algorithm:

1. Hitung in-degree untuk setiap node.
2. Masukkan node dengan in-degree 0 ke queue.
3. Proses queue: kurangi in-degree tetangga.
4. Jika in-degree tetangga menjadi 0, tambahkan ke queue.
5. Jika jumlah node terurut != jumlah total node, terdeteksi siklus.

## 8. Usage Example

```go
g := graph.New()

// Add nodes
g.AddNode(graph.Node{ID: "spec", Kind: graph.NodeKindModule, Name: "spec-parser"})
g.AddNode(graph.Node{ID: "norm", Kind: graph.NodeKindModule, Name: "normalizer"})
g.AddNode(graph.Node{ID: "gen", Kind: graph.NodeKindModule, Name: "generator"})

// Add edges (dependencies)
g.AddEdge(graph.Edge{From: "spec", To: "norm", Kind: graph.EdgeKindDependency})
g.AddEdge(graph.Edge{From: "norm", To: "gen", Kind: graph.EdgeKindDependency})

// Topological sort
sorted, err := g.TopologicalSort()
// [spec, norm, gen]
```

## 9. Acceptance Criteria
- A system dependency can be analyzed through the graph representation.
- Topological sort produces correct execution order.
- Cycles are detected and reported as errors.
- The graph can be used by planning and validation components without manual conversion.

package graph

import (
	"testing"
)

func TestNewGraphIsEmpty(t *testing.T) {
	g := New()
	if g.NodeCount() != 0 {
		t.Fatalf("expected 0 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Fatalf("expected 0 edges, got %d", g.EdgeCount())
	}
}

func TestAddNode(t *testing.T) {
	g := New()
	err := g.AddNode(Node{ID: "a", Kind: NodeKindModule, Name: "module-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.NodeCount() != 1 {
		t.Fatalf("expected 1 node, got %d", g.NodeCount())
	}
}

func TestAddNodeEmptyID(t *testing.T) {
	g := New()
	err := g.AddNode(Node{ID: "", Kind: NodeKindModule})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestAddNodeDuplicate(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	err := g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	if err == nil {
		t.Fatal("expected error for duplicate node")
	}
}

func TestAddEdge(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	err := g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.EdgeCount() != 1 {
		t.Fatalf("expected 1 edge, got %d", g.EdgeCount())
	}
}

func TestAddEdgeMissingNode(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	err := g.AddEdge(Edge{From: "a", To: "missing", Kind: EdgeKindDependency})
	if err == nil {
		t.Fatal("expected error for missing target node")
	}
}

func TestAddEdgeDuplicate(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	err := g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	if err == nil {
		t.Fatal("expected error for duplicate edge")
	}
}

func TestRemoveNode(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	err := g.RemoveNode("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.NodeCount() != 1 {
		t.Fatalf("expected 1 node after removal, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Fatalf("expected 0 edges after node removal, got %d", g.EdgeCount())
	}
}

func TestGetNeighbors(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "c", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	_ = g.AddEdge(Edge{From: "a", To: "c", Kind: EdgeKindDependency})
	neighbors := g.GetNeighbors("a")
	if len(neighbors) != 2 {
		t.Fatalf("expected 2 neighbors, got %d", len(neighbors))
	}
}

func TestTopologicalSort(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "c", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	_ = g.AddEdge(Edge{From: "b", To: "c", Kind: EdgeKindDependency})

	sorted, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sorted) != 3 {
		t.Fatalf("expected 3 nodes in sort, got %d", len(sorted))
	}

	index := make(map[string]int)
	for i, n := range sorted {
		index[n.ID] = i
	}
	if index["a"] >= index["b"] || index["b"] >= index["c"] {
		t.Fatalf("topological order incorrect: %v", sorted)
	}
}

func TestTopologicalSortDetectsCycle(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	_ = g.AddEdge(Edge{From: "b", To: "a", Kind: EdgeKindDependency})

	_, err := g.TopologicalSort()
	if err == nil {
		t.Fatal("expected cycle detection error")
	}
}

func TestHasCycle(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})

	if g.HasCycle() {
		t.Fatal("expected no cycle")
	}

	_ = g.AddEdge(Edge{From: "b", To: "a", Kind: EdgeKindDependency})
	if !g.HasCycle() {
		t.Fatal("expected cycle detected")
	}
}

func TestInDegree(t *testing.T) {
	g := New()
	_ = g.AddNode(Node{ID: "a", Kind: NodeKindModule})
	_ = g.AddNode(Node{ID: "b", Kind: NodeKindModule})
	_ = g.AddEdge(Edge{From: "a", To: "b", Kind: EdgeKindDependency})
	_ = g.AddEdge(Edge{From: "b", To: "a", Kind: EdgeKindDependency})

	if g.InDegree("a") != 1 {
		t.Fatalf("expected in-degree 1 for a, got %d", g.InDegree("a"))
	}
	if g.OutDegree("a") != 1 {
		t.Fatalf("expected out-degree 1 for a, got %d", g.OutDegree("a"))
	}
}

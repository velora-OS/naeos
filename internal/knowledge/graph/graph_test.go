package graph

import (
	"testing"
)

func TestNewKnowledgeGraph(t *testing.T) {
	kg := New()
	if kg == nil {
		t.Fatal("expected non-nil knowledge graph")
	}
	if kg.NodeCount() != 0 {
		t.Fatalf("expected 0 nodes, got %d", kg.NodeCount())
	}
}

func TestAddNode(t *testing.T) {
	kg := New()
	err := kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision, Topic: "architecture", Component: "api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kg.NodeCount() != 1 {
		t.Fatalf("expected 1 node, got %d", kg.NodeCount())
	}
}

func TestAddNodeEmptyID(t *testing.T) {
	kg := New()
	err := kg.AddNode(Node{ID: "", Type: NodeTypeDecision})
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestAddNodeDuplicate(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	err := kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	if err == nil {
		t.Fatal("expected error for duplicate node")
	}
}

func TestGetNode(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision, Topic: "test"})
	node, ok := kg.GetNode("n1")
	if !ok {
		t.Fatal("expected to find node")
	}
	if node.Topic != "test" {
		t.Fatalf("expected topic 'test', got %s", node.Topic)
	}
}

func TestRemoveNode(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})

	err := kg.RemoveNode("n1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kg.NodeCount() != 1 {
		t.Fatalf("expected 1 node, got %d", kg.NodeCount())
	}
	if kg.EdgeCount() != 0 {
		t.Fatalf("expected 0 edges, got %d", kg.EdgeCount())
	}
}

func TestAddEdge(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	err := kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kg.EdgeCount() != 1 {
		t.Fatalf("expected 1 edge, got %d", kg.EdgeCount())
	}
}

func TestAddEdgeMissingNode(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	err := kg.AddEdge(Edge{From: "n1", To: "missing", Type: EdgeTypeDependsOn})
	if err == nil {
		t.Fatal("expected error for missing target node")
	}
}

func TestFindByTopic(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision, Topic: "auth"})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeDecision, Topic: "api"})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeDecision, Topic: "auth"})

	nodes := kg.FindByTopic("auth")
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestFindByComponent(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeComponent, Component: "api"})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent, Component: "db"})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeComponent, Component: "api"})

	nodes := kg.FindByComponent("api")
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestFindByType(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeDecision})

	nodes := kg.FindByType(NodeTypeDecision)
	if len(nodes) != 2 {
		t.Fatalf("expected 2 decision nodes, got %d", len(nodes))
	}
}

func TestGetNeighbors(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})
	_ = kg.AddEdge(Edge{From: "n1", To: "n3", Type: EdgeTypeDependsOn})

	neighbors := kg.GetNeighbors("n1")
	if len(neighbors) != 2 {
		t.Fatalf("expected 2 neighbors, got %d", len(neighbors))
	}
}

func TestNodesAndEdges(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})

	nodes := kg.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}

	edges := kg.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
}

func TestGetPredecessors(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})

	predecessors := kg.GetPredecessors("n2")
	if len(predecessors) != 1 {
		t.Fatalf("expected 1 predecessor, got %d", len(predecessors))
	}
	if predecessors[0].ID != "n1" {
		t.Fatalf("expected predecessor n1, got %s", predecessors[0].ID)
	}
}

func TestGetConnected(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})
	_ = kg.AddEdge(Edge{From: "n2", To: "n3", Type: EdgeTypeDependsOn})

	connected := kg.GetConnected("n1")
	if len(connected) != 1 {
		t.Fatalf("expected 1 connected node (n2), got %d", len(connected))
	}
}

func TestFindByEdgeType(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})
	_ = kg.AddEdge(Edge{From: "n1", To: "n3", Type: EdgeTypeImplements})

	edges := kg.FindByEdgeType(EdgeTypeDependsOn)
	if len(edges) != 1 {
		t.Fatalf("expected 1 depends_on edge, got %d", len(edges))
	}
}

func TestFindEdgesFrom(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})

	edges := kg.FindEdgesFrom("n1")
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge from n1, got %d", len(edges))
	}
}

func TestFindEdgesTo(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})

	edges := kg.FindEdgesTo("n2")
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge to n2, got %d", len(edges))
	}
}

func TestFindByMetadata(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision, Metadata: map[string]string{"env": "prod"}})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeDecision, Metadata: map[string]string{"env": "staging"}})

	nodes := kg.FindByMetadata("env", "prod")
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
}

func TestHasPath(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeDecision})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeComponent})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeComponent})
	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeDependsOn})
	_ = kg.AddEdge(Edge{From: "n2", To: "n3", Type: EdgeTypeDependsOn})

	if !kg.HasPath("n1", "n3") {
		t.Fatal("expected path from n1 to n3")
	}
	if kg.HasPath("n3", "n1") {
		t.Fatal("expected no path from n3 to n1")
	}
}

func TestNewNodeTypes(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "svc1", Type: NodeTypeService, Topic: "api"})
	_ = kg.AddNode(Node{ID: "mod1", Type: NodeTypeModule, Topic: "auth"})
	_ = kg.AddNode(Node{ID: "api1", Type: NodeTypeAPI, Topic: "rest"})
	_ = kg.AddNode(Node{ID: "db1", Type: NodeTypeStorage, Topic: "postgres"})
	_ = kg.AddNode(Node{ID: "dep1", Type: NodeTypeDeployment, Topic: "k8s"})
	_ = kg.AddNode(Node{ID: "test1", Type: NodeTypeTesting, Topic: "unit"})
	_ = kg.AddNode(Node{ID: "sec1", Type: NodeTypeSecurity, Topic: "auth"})

	if kg.NodeCount() != 7 {
		t.Fatalf("expected 7 nodes, got %d", kg.NodeCount())
	}

	services := kg.FindByType(NodeTypeService)
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}
}

func TestNewEdgeTypes(t *testing.T) {
	kg := New()
	_ = kg.AddNode(Node{ID: "n1", Type: NodeTypeService})
	_ = kg.AddNode(Node{ID: "n2", Type: NodeTypeModule})
	_ = kg.AddNode(Node{ID: "n3", Type: NodeTypeStorage})

	_ = kg.AddEdge(Edge{From: "n1", To: "n2", Type: EdgeTypeContains})
	_ = kg.AddEdge(Edge{From: "n1", To: "n3", Type: EdgeTypeConnectsTo})

	contains := kg.FindByEdgeType(EdgeTypeContains)
	if len(contains) != 1 {
		t.Fatalf("expected 1 contains edge, got %d", len(contains))
	}
}

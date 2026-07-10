package graph

import (
	"fmt"
	"sync"
)

type NodeType string

const (
	NodeTypeDecision      NodeType = "decision"
	NodeTypeRequirement   NodeType = "requirement"
	NodeTypeRationale     NodeType = "rationale"
	NodeTypeComponent     NodeType = "component"
	NodeTypePolicy        NodeType = "policy"
	NodeTypeImplementation NodeType = "implementation"
	NodeTypeHistorical    NodeType = "historical"
	NodeTypeService       NodeType = "service"
	NodeTypeModule        NodeType = "module"
	NodeTypeAPI           NodeType = "api"
	NodeTypeStorage       NodeType = "storage"
	NodeTypeDeployment    NodeType = "deployment"
	NodeTypeTesting       NodeType = "testing"
	NodeTypeSecurity      NodeType = "security"
)

type EdgeType string

const (
	EdgeTypeDependsOn     EdgeType = "depends_on"
	EdgeTypeImplements    EdgeType = "implements"
	EdgeTypeRelatedTo     EdgeType = "related_to"
	EdgeTypeSupersedes    EdgeType = "supersedes"
	EdgeTypeConflictsWith EdgeType = "conflicts_with"
	EdgeTypeContains      EdgeType = "contains"
	EdgeTypeExposes       EdgeType = "exposes"
	EdgeTypeConnectsTo    EdgeType = "connects_to"
	EdgeTypeDeploysTo     EdgeType = "deploys_to"
	EdgeTypeTests         EdgeType = "tests"
	EdgeTypeSecures       EdgeType = "secures"
	EdgeTypeUses          EdgeType = "uses"
	EdgeTypeExtends       EdgeType = "extends"
)

type Node struct {
	ID        string
	Type      NodeType
	Topic     string
	Component string
	Version   string
	Content   string
	Metadata  map[string]string
}

type Edge struct {
	From   string
	To     string
	Type   EdgeType
	Weight int
}

type KnowledgeGraph struct {
	mu    sync.RWMutex
	nodes map[string]*Node
	edges []Edge
}

func New() *KnowledgeGraph {
	return &KnowledgeGraph{
		nodes: make(map[string]*Node),
	}
}

func (kg *KnowledgeGraph) AddNode(n Node) error {
	kg.mu.Lock()
	defer kg.mu.Unlock()

	if n.ID == "" {
		return fmt.Errorf("node ID must not be empty")
	}
	if _, exists := kg.nodes[n.ID]; exists {
		return fmt.Errorf("node %s already exists", n.ID)
	}
	kg.nodes[n.ID] = &n
	return nil
}

func (kg *KnowledgeGraph) GetNode(id string) (*Node, bool) {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	n, ok := kg.nodes[id]
	return n, ok
}

func (kg *KnowledgeGraph) RemoveNode(id string) error {
	kg.mu.Lock()
	defer kg.mu.Unlock()

	if _, exists := kg.nodes[id]; !exists {
		return fmt.Errorf("node %s not found", id)
	}
	delete(kg.nodes, id)

	filtered := kg.edges[:0]
	for _, e := range kg.edges {
		if e.From != id && e.To != id {
			filtered = append(filtered, e)
		}
	}
	kg.edges = filtered
	return nil
}

func (kg *KnowledgeGraph) AddEdge(e Edge) error {
	kg.mu.Lock()
	defer kg.mu.Unlock()

	if _, exists := kg.nodes[e.From]; !exists {
		return fmt.Errorf("source node %s not found", e.From)
	}
	if _, exists := kg.nodes[e.To]; !exists {
		return fmt.Errorf("target node %s not found", e.To)
	}
	kg.edges = append(kg.edges, e)
	return nil
}

func (kg *KnowledgeGraph) FindByTopic(topic string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if n.Topic == topic {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByComponent(component string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if n.Component == component {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByVersion(version string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if n.Version == version {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByType(nodeType NodeType) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if n.Type == nodeType {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByEdgeType(edgeType EdgeType) []Edge {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []Edge
	for _, e := range kg.edges {
		if e.Type == edgeType {
			result = append(result, e)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindEdgesFrom(id string) []Edge {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []Edge
	for _, e := range kg.edges {
		if e.From == id {
			result = append(result, e)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindEdgesTo(id string) []Edge {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []Edge
	for _, e := range kg.edges {
		if e.To == id {
			result = append(result, e)
		}
	}
	return result
}

func (kg *KnowledgeGraph) Nodes() []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	result := make([]*Node, 0, len(kg.nodes))
	for _, n := range kg.nodes {
		result = append(result, n)
	}
	return result
}

func (kg *KnowledgeGraph) Edges() []Edge {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	result := make([]Edge, len(kg.edges))
	copy(result, kg.edges)
	return result
}

func (kg *KnowledgeGraph) NodeCount() int {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	return len(kg.nodes)
}

func (kg *KnowledgeGraph) EdgeCount() int {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	return len(kg.edges)
}

func (kg *KnowledgeGraph) GetNeighbors(id string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, e := range kg.edges {
		if e.From == id {
			if n, ok := kg.nodes[e.To]; ok {
				result = append(result, n)
			}
		}
	}
	return result
}

func (kg *KnowledgeGraph) GetPredecessors(id string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, e := range kg.edges {
		if e.To == id {
			if n, ok := kg.nodes[e.From]; ok {
				result = append(result, n)
			}
		}
	}
	return result
}

func (kg *KnowledgeGraph) GetConnected(id string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	seen := make(map[string]bool)
	var result []*Node
	for _, e := range kg.edges {
		if e.From == id {
			if !seen[e.To] {
				if n, ok := kg.nodes[e.To]; ok {
					result = append(result, n)
					seen[e.To] = true
				}
			}
		}
		if e.To == id {
			if !seen[e.From] {
				if n, ok := kg.nodes[e.From]; ok {
					result = append(result, n)
					seen[e.From] = true
				}
			}
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByMetadata(key, value string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if n.Metadata != nil && n.Metadata[key] == value {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) FindByContentSubstring(substring string) []*Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	var result []*Node
	for _, n := range kg.nodes {
		if contains := false; contains {
			result = append(result, n)
		}
	}
	return result
}

func (kg *KnowledgeGraph) HasPath(from, to string) bool {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	visited := make(map[string]bool)
	queue := []string{from}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == to {
			return true
		}

		for _, e := range kg.edges {
			if e.From == current && !visited[e.To] {
				visited[e.To] = true
				queue = append(queue, e.To)
			}
		}
	}
	return false
}

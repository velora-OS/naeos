package graph

import "fmt"

type NodeKind string

const (
	NodeKindService      NodeKind = "service"
	NodeKindModule       NodeKind = "module"
	NodeKindComponent    NodeKind = "component"
	NodeKindAPI          NodeKind = "api"
	NodeKindStorage      NodeKind = "storage"
	NodeKindInfra        NodeKind = "infrastructure"
	NodeKindDocumentation NodeKind = "documentation"
	NodeKindTesting      NodeKind = "testing"
	NodeKindDeployment   NodeKind = "deployment"
)

type EdgeKind string

const (
	EdgeKindDependency  EdgeKind = "dependency"
	EdgeKindExecution   EdgeKind = "execution"
	EdgeKindDataflow    EdgeKind = "dataflow"
	EdgeKindPolicy      EdgeKind = "policy"
)

type Node struct {
	ID   string
	Kind NodeKind
	Name string
}

type Edge struct {
	From string
	To   string
	Kind EdgeKind
}

type PlannerGraph struct {
	nodes map[string]Node
	edges []Edge
}

func New() *PlannerGraph {
	return &PlannerGraph{
		nodes: make(map[string]Node),
	}
}

func (g *PlannerGraph) AddNode(n Node) error {
	if n.ID == "" {
		return fmt.Errorf("node ID must not be empty")
	}
	if _, exists := g.nodes[n.ID]; exists {
		return fmt.Errorf("node %s already exists", n.ID)
	}
	g.nodes[n.ID] = n
	return nil
}

func (g *PlannerGraph) GetNode(id string) (Node, bool) {
	n, ok := g.nodes[id]
	return n, ok
}

func (g *PlannerGraph) RemoveNode(id string) error {
	if _, exists := g.nodes[id]; !exists {
		return fmt.Errorf("node %s not found", id)
	}
	delete(g.nodes, id)
	filtered := g.edges[:0]
	for _, e := range g.edges {
		if e.From != id && e.To != id {
			filtered = append(filtered, e)
		}
	}
	g.edges = filtered
	return nil
}

func (g *PlannerGraph) AddEdge(e Edge) error {
	if _, exists := g.nodes[e.From]; !exists {
		return fmt.Errorf("source node %s not found", e.From)
	}
	if _, exists := g.nodes[e.To]; !exists {
		return fmt.Errorf("target node %s not found", e.To)
	}
	for _, existing := range g.edges {
		if existing.From == e.From && existing.To == e.To && existing.Kind == e.Kind {
			return fmt.Errorf("edge from %s to %s with kind %s already exists", e.From, e.To, e.Kind)
		}
	}
	g.edges = append(g.edges, e)
	return nil
}

func (g *PlannerGraph) GetNeighbors(id string) []Node {
	var neighbors []Node
	for _, e := range g.edges {
		if e.From == id {
			if n, ok := g.nodes[e.To]; ok {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}

func (g *PlannerGraph) Nodes() []Node {
	result := make([]Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		result = append(result, n)
	}
	return result
}

func (g *PlannerGraph) Edges() []Edge {
	result := make([]Edge, len(g.edges))
	copy(result, g.edges)
	return result
}

func (g *PlannerGraph) NodeCount() int {
	return len(g.nodes)
}

func (g *PlannerGraph) EdgeCount() int {
	return len(g.edges)
}

func (g *PlannerGraph) TopologicalSort() ([]Node, error) {
	inDegree := make(map[string]int)
	for id := range g.nodes {
		inDegree[id] = 0
	}
	for _, e := range g.edges {
		inDegree[e.To]++
	}

	var queue []string
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []Node
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		sorted = append(sorted, g.nodes[id])

		for _, e := range g.edges {
			if e.From == id {
				inDegree[e.To]--
				if inDegree[e.To] == 0 {
					queue = append(queue, e.To)
				}
			}
		}
	}

	if len(sorted) != len(g.nodes) {
		return nil, fmt.Errorf("cycle detected in graph")
	}

	return sorted, nil
}

func (g *PlannerGraph) HasCycle() bool {
	_, err := g.TopologicalSort()
	return err != nil
}

func (g *PlannerGraph) InDegree(id string) int {
	count := 0
	for _, e := range g.edges {
		if e.To == id {
			count++
		}
	}
	return count
}

func (g *PlannerGraph) OutDegree(id string) int {
	count := 0
	for _, e := range g.edges {
		if e.From == id {
			count++
		}
	}
	return count
}

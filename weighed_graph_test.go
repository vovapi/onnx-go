package onnx

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

const (
	self, absent = math.MaxFloat64, float64(-1)
)

type edge struct {
	from   *nodeTest
	to     *nodeTest
	weight float64
}

func newExpectedGraph(e []edge) *testExpectedGraph {
	g := &testExpectedGraph{
		nodes: make(map[int64]*nodeTest),
		from:  make(map[int64]map[int64]graph.WeightedEdge),
		to:    make(map[int64]map[int64]graph.WeightedEdge),
	}
	for _, e := range e {
		if _, ok := g.nodes[e.from.id]; !ok {
			g.nodes[e.from.id] = e.from
		}
		if _, ok := g.nodes[e.to.id]; !ok {
			g.nodes[e.to.id] = e.to
		}
		g.SetWeightedEdge(g.NewWeightedEdge(e.from, e.to, e.weight))

	}
	return g
}

// testExpectedGraph implements a generalized weighted directed graph.
type testExpectedGraph struct {
	nodes map[int64]*nodeTest
	from  map[int64]map[int64]graph.WeightedEdge
	to    map[int64]map[int64]graph.WeightedEdge
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *testExpectedGraph) Edge(uid, vid int64) graph.Edge {
	return g.WeightedEdge(uid, vid)
}

// Edges returns all the edges in the graph.
func (g *testExpectedGraph) Edges() graph.Edges {
	var edges []graph.Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedEdges(edges)
}

// From returns all nodes in g that can be reached directly from n.
func (g *testExpectedGraph) From(id int64) graph.Nodes {
	if _, ok := g.from[id]; !ok {
		return graph.Empty
	}

	from := make([]graph.Node, len(g.from[id]))
	i := 0
	for vid := range g.from[id] {
		from[i] = g.nodes[vid]
		i++
	}
	if len(from) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(from)
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (g *testExpectedGraph) HasEdgeBetween(xid, yid int64) bool {
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (g *testExpectedGraph) HasEdgeFromTo(uid, vid int64) bool {
	if _, ok := g.from[uid][vid]; !ok {
		return false
	}
	return true
}

// NewWeightedEdge returns a new weighted edge from the source to the destination node.
func (g *testExpectedGraph) NewWeightedEdge(from, to graph.Node, weight float64) graph.WeightedEdge {
	return &simple.WeightedEdge{F: from, T: to, W: weight}
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *testExpectedGraph) Node(id int64) graph.Node {
	return g.nodes[id]
}

// Nodes returns all the nodes in the graph.
func (g *testExpectedGraph) Nodes() graph.Nodes {
	if len(g.from) == 0 {
		return graph.Empty
	}
	nodes := make([]graph.Node, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}
	return iterator.NewOrderedNodes(nodes)
}

// To returns all nodes in g that can reach directly to n.
func (g *testExpectedGraph) To(id int64) graph.Nodes {
	if _, ok := g.from[id]; !ok {
		return graph.Empty
	}

	to := make([]graph.Node, len(g.to[id]))
	i := 0
	for uid := range g.to[id] {
		to[i] = g.nodes[uid]
		i++
	}
	if len(to) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedNodes(to)
}

// Weight returns the weight for the edge between x and y if Edge(x, y) returns a non-nil Edge.
// If x and y are the same node or there is no joining edge between the two nodes the weight
// value returned is either the graph's absent or self value. Weight returns true if an edge
// exists between x and y or if x and y have the same ID, false otherwise.
func (g *testExpectedGraph) Weight(xid, yid int64) (w float64, ok bool) {
	if xid == yid {
		return self, true
	}
	if to, ok := g.from[xid]; ok {
		if e, ok := to[yid]; ok {
			return e.Weight(), true
		}
	}
	return absent, false
}

// WeightedEdge returns the weighted edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *testExpectedGraph) WeightedEdge(uid, vid int64) graph.WeightedEdge {
	edge, ok := g.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

// WeightedEdges returns all the weighted edges in the graph.
func (g *testExpectedGraph) WeightedEdges() graph.WeightedEdges {
	var edges []graph.WeightedEdge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedWeightedEdges(edges)
}

// SetWeightedEdge adds a weighted edge from one node to another. If the nodes do not exist, they are added
// and are set to the nodes of the edge otherwise.
// It will panic if the IDs of the e.From and e.To are equal.
func (g *testExpectedGraph) SetWeightedEdge(e graph.WeightedEdge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if g.from[fid] == nil {
		g.from[fid] = make(map[int64]graph.WeightedEdge)
	}
	if g.to[tid] == nil {
		g.to[tid] = make(map[int64]graph.WeightedEdge)
	}
	g.from[fid][tid] = e
	g.to[tid][fid] = e
}

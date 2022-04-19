package adjacencylist

import "fmt"

type (
	VertexDescriptor int32
	EdgeDestriptor   int32
	GraphOptions     func(*GraphAttribute)
	AttributeNone    struct{}
	EdgePair         struct {
		from VertexDescriptor
		to   VertexDescriptor
	}
)

type GraphAttribute struct {
	isDirected bool
}

func DirectedGraph(ga *GraphAttribute) {
	ga.isDirected = true
}

type Graph[VertexAttribute any, EdgeAttribute any] struct {
	verticesAttributes map[VertexDescriptor]*VertexAttribute
	edgesAttributes    map[EdgeDestriptor]*EdgeAttribute
	verticesAdgacency  map[VertexDescriptor][]EdgeDestriptor
	edgeAdjacency      map[EdgeDestriptor]EdgePair
	attributes         *GraphAttribute
}

func NewGraph[VA any, EA any](args ...GraphOptions) *Graph[VA, EA] {
	ga := &GraphAttribute{
		isDirected: false,
	}
	for _, arg := range args {
		arg(ga)
	}
	return &Graph[VA, EA]{
		attributes:         ga,
		verticesAttributes: make(map[VertexDescriptor]*VA),
		edgesAttributes:    make(map[EdgeDestriptor]*EA),
		verticesAdgacency:  make(map[VertexDescriptor][]EdgeDestriptor),
		edgeAdjacency:      make(map[EdgeDestriptor]EdgePair),
	}
}

func (g *Graph[VA, EA]) AddVertex(attributes ...*VA) VertexDescriptor {
	var attr *VA
	if len(attributes) > 0 {
		attr = attributes[0]
	} else {
		attr = new(VA)
	}
	v := VertexDescriptor(len(g.verticesAdgacency))
	g.verticesAdgacency[v] = []EdgeDestriptor{}
	g.verticesAttributes[v] = attr
	return v
}

func (g *Graph[VA, EA]) AddEdge(from, to VertexDescriptor, attributes ...*EA) EdgeDestriptor {
	var attr *EA
	if len(attributes) > 0 {
		attr = attributes[0]
	} else {
		attr = new(EA)
	}
	e := EdgeDestriptor(len(g.edgeAdjacency))
	g.verticesAdgacency[from] = append(g.verticesAdgacency[from], e)
	g.edgeAdjacency[e] = EdgePair{
		from: from,
		to:   to,
	}
	if !g.attributes.isDirected {
		g.verticesAdgacency[to] = append(g.verticesAdgacency[to], e)
	}
	g.edgesAttributes[e] = attr
	return e
}

func (g *Graph[VA, EA]) GetVertices() []VertexDescriptor {
	var vertices []VertexDescriptor
	for v := range g.verticesAdgacency {
		vertices = append(vertices, v)
	}
	return vertices
}

func (g *Graph[VA, EA]) GetVertexAttribute(v VertexDescriptor) *VA {
	return g.verticesAttributes[v]
}

func (g *Graph[VA, EA]) GetEdgeAttribute(e EdgeDestriptor) *EA {
	return g.edgesAttributes[e]
}

type LabelWritter[A any] func(*A) string

func DefaultLabelWritter[A any](*A) string {
	return ""
}

type LabelWritterOption[VA, EA any] func(*LabelWritters[VA, EA])

type LabelWritters[VA, EA any] struct {
	vertexWritter LabelWritter[VA]
	edgeWritter   LabelWritter[EA]
}

func WithVertexLabelWritter[VA, EA any](lw LabelWritter[VA]) LabelWritterOption[VA, EA] {
	return func(lws *LabelWritters[VA, EA]) {
		lws.vertexWritter = lw
	}
}

func WithEdgeLabelWritter[VA, EA any](lw LabelWritter[EA]) LabelWritterOption[VA, EA] {
	return func(lws *LabelWritters[VA, EA]) {
		lws.edgeWritter = lw
	}
}

func (g *Graph[VA, EA]) DumpGraphviz(opts ...LabelWritterOption[VA, EA]) string {
	var dot string
	var arrow string
	if g.attributes.isDirected {
		dot += "digraph g {\n"
		arrow = "->"
	} else {
		dot += "graph g {\n"
		arrow = "--"
	}

	lws := &LabelWritters[VA, EA]{
		vertexWritter: DefaultLabelWritter[VA],
		edgeWritter:   DefaultLabelWritter[EA],
	}
	for _, opt := range opts {
		opt(lws)
	}

	for i := 0; i < len(g.verticesAdgacency); i++ {
		dot += fmt.Sprintf("%d%s;\n", i,
			lws.vertexWritter(g.verticesAttributes[VertexDescriptor(i)]))
	}

	for ei, e := range g.edgeAdjacency {
		dot += fmt.Sprintf("%d %s %d%s;\n", e.from, arrow, e.to,
			lws.edgeWritter(g.edgesAttributes[EdgeDestriptor(ei)]))
	}

	dot += "}\n"
	return dot
}

func (g *Graph[VA, EV]) Target(e EdgeDestriptor) (VertexDescriptor, error) {
	if a, ok := g.edgeAdjacency[e]; ok {
		return a.to, nil
	}
	return VertexDescriptor(0), fmt.Errorf("Could not find edge %d", e)
}

func (g *Graph[VA, EV]) Source(e EdgeDestriptor) (VertexDescriptor, error) {
	if a, ok := g.edgeAdjacency[e]; ok {
		return a.from, nil
	}
	return VertexDescriptor(0), fmt.Errorf("Could not find edge %d", e)
}

func (g *Graph[VA, EV]) OutEdges(v VertexDescriptor) ([]EdgeDestriptor, error) {
	if a, ok := g.verticesAdgacency[v]; ok {
		return a, nil
	}
	return []EdgeDestriptor{}, fmt.Errorf("Could not find vertex %d", v)
}

func (g *Graph[VA, EV]) IsLeef(v VertexDescriptor) bool {
	if a, ok := g.verticesAdgacency[v]; ok {
		return len(a) == 0
	}
	return false
}

func (g *Graph[VA, EV]) Neighbors(v VertexDescriptor) ([]VertexDescriptor, error) {
	oe, err := g.OutEdges(v)
	if err != nil {
		return nil, err
	}
	var neighbors []VertexDescriptor
	for _, ed := range oe {
		target, err := g.Target(ed)
		if err != nil {
			return nil, err
		}
		neighbors = append(neighbors, target)
	}
	return neighbors, nil
}

func (g *Graph[VA, EV]) HasCycle() (bool, []VertexDescriptor, error) {
	visited := make(map[VertexDescriptor]bool)
	recStack := make(map[VertexDescriptor]bool)
	for v := range g.verticesAdgacency {
		visited[v] = false
		recStack[v] = false
	}

	var isCyclicRec func(v VertexDescriptor, stack []VertexDescriptor) (bool, []VertexDescriptor, error)
	isCyclicRec = func(v VertexDescriptor, stack []VertexDescriptor) (bool, []VertexDescriptor, error) {
		visited[v] = true
		recStack[v] = true

		neighbors, err := g.Neighbors(v)
		if err != nil {
			return false, nil, err
		}
		for _, neighbor := range neighbors {
			if recStack[neighbor] {
				return true, append(stack, v), nil
			}
			if !visited[neighbor] {
				if isC, nstack, err := isCyclicRec(neighbor, append(stack, v)); err != nil {
					return false, nil, err
				} else if isC {
					return true, nstack, nil
				}
			}
		}
		recStack[v] = false

		return false, nil, nil
	}

	for v := range g.verticesAdgacency {
		if !visited[v] {
			c, stack, err := isCyclicRec(v, []VertexDescriptor{})
			if err != nil {
				return false, nil, err
			}
			if c {
				return true, stack, nil
			}
		}
	}
	return false, nil, nil
}

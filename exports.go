package graph

import "errors"

type Color uint8

const (
	ColorWhite Color = iota
	ColorGray
	ColorBlack
)

var ErrUndefinedNode = errors.New("undefined node")
var ErrEdgeConflicts = errors.New("edge conflicts")

type (
	ID     string
	Weight uint64
)

type Graph interface {
	Nodes() []ID
	NodesLen() int
	CheckNode(id ID) bool
	AddNode(id ID) bool
	DeleteNode(id ID) bool
	ChildrenOf(id ID) ([]ID, bool)
	ParentsOf(id ID) ([]ID, bool)
	GetEdge(src, tgt ID) (Weight, bool)
	AddEdge(src, tgt ID, w Weight) error
	SetEdge(src, tgt ID, w Weight) error
	DeleteEdge(src, tgt ID) bool
}

type Runtime struct {
	graph Graph
}

func NewRuntime(graph Graph) *Runtime {
	return &Runtime{graph: graph}
}

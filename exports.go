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

type Weight uint64

type Graph[ID comparable] interface {
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

type Runtime[ID comparable] struct {
	graph Graph[ID]
}

func NewRuntime[ID comparable](graph Graph[ID]) *Runtime[ID] {
	return &Runtime[ID]{graph: graph}
}

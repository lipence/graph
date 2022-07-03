package graph

import (
	"sync"

	"github.com/tidwall/hashmap"
)

type baseGraph[ID comparable] struct {
	lock  sync.RWMutex
	nodes hashmap.Set[ID]
	edges hashmap.Map[[2]ID, Weight]
}

func (g *baseGraph[ID]) exportNodes() []ID { return g.nodes.Keys() }

func (g *baseGraph[ID]) countNodes() int { return g.nodes.Len() }

func (g *baseGraph[ID]) checkNode(id ID) bool {
	return g.nodes.Contains(id)
}

func (g *baseGraph[ID]) addNode(id ID) bool {
	if g.nodes.Contains(id) {
		return false
	}
	g.nodes.Insert(id)
	return true
}

func (g *baseGraph[ID]) deleteNode(id ID) bool {
	if !g.nodes.Contains(id) {
		return false
	} else {
		g.nodes.Delete(id)
		return true
	}
}

func (g *baseGraph[ID]) getEdge(src, tgt ID) (Weight, bool) {
	return g.edges.Get([2]ID{src, tgt})
}

func (g *baseGraph[ID]) addEdge(src, tgt ID, w Weight) bool {
	var edgeKey = [2]ID{src, tgt}
	if isOK(g.edges.Get(edgeKey)) {
		return false
	}
	g.edges.Set(edgeKey, w)
	return true
}

func (g *baseGraph[ID]) setEdge(src, tgt ID, w Weight) {
	g.edges.Set([2]ID{src, tgt}, w)
}

func (g *baseGraph[ID]) deleteEdge(src, tgt ID) bool {
	var edgeKey = [2]ID{src, tgt}
	if !isOK(g.edges.Get(edgeKey)) {
		return false
	}
	g.edges.Delete(edgeKey)
	return true
}

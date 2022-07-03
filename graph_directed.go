package graph

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/tidwall/hashmap"
)

type directedGraph[ID comparable] struct {
	baseGraph[ID]
	opexCount    uint64
	cacheLock    sync.Mutex
	cachedOpex   uint64
	cacheBySrcID hashmap.Map[ID, *hashmap.Set[ID]]
	cacheByTgtID hashmap.Map[ID, *hashmap.Set[ID]]
}

func NewDirected[ID comparable]() *directedGraph[ID] { return &directedGraph[ID]{} }

func (g *directedGraph[ID]) checkEdgeNodes(src, tgt ID) bool {
	return g.baseGraph.checkNode(src) && g.baseGraph.checkNode(tgt)
}

func (g *directedGraph[ID]) checkEdgeNodesWithErr(src, tgt ID) error {
	if !g.baseGraph.checkNode(src) {
		return fmt.Errorf("%w: id = %v", ErrUndefinedNode, src)
	}
	if !g.baseGraph.checkNode(tgt) {
		return fmt.Errorf("%w: id = %v", ErrUndefinedNode, tgt)
	}
	return nil
}

func (g *directedGraph[ID]) guaranteeCached() (cm, pm hashmap.Map[ID, *hashmap.Set[ID]]) {
	g.cacheLock.Lock()
	defer g.cacheLock.Unlock()
	if currentOpex := g.opexCount; g.cachedOpex == currentOpex {
		return g.cacheBySrcID, g.cacheByTgtID
	} else {
		defer func() { g.cachedOpex = currentOpex }()
	}
	cm, pm = hashmap.Map[ID, *hashmap.Set[ID]]{}, hashmap.Map[ID, *hashmap.Set[ID]]{}
	g.baseGraph.edges.Scan(func(ids [2]ID, _ Weight) bool {
		var exist bool
		var cs, ps *hashmap.Set[ID]
		// children
		cs, exist = cm.Get(ids[0])
		if !exist {
			cs = &hashmap.Set[ID]{}
			cm.Set(ids[0], cs)
		}
		cs.Insert(ids[1])
		// parent
		ps, exist = pm.Get(ids[1])
		if !exist {
			ps = &hashmap.Set[ID]{}
			pm.Set(ids[1], ps)
		}
		ps.Insert(ids[0])
		return true
	})
	g.cacheBySrcID, g.cacheByTgtID = cm, pm
	return cm, pm
}

func (g *directedGraph[ID]) exportCache() (children, parents map[ID][]ID) {
	var cm, pm = g.guaranteeCached()
	children, parents = map[ID][]ID{}, map[ID][]ID{}
	cm.Scan(func(i ID, s *hashmap.Set[ID]) bool {
		children[i] = s.Keys()
		return true
	})
	pm.Scan(func(i ID, s *hashmap.Set[ID]) bool {
		parents[i] = s.Keys()
		return true
	})
	return children, parents
}

func (g *directedGraph[ID]) ChildrenOf(id ID) ([]ID, bool) {
	if !g.checkNode(id) {
		return nil, false
	}
	var cm, _ = g.guaranteeCached()
	if childrenSet, ok := cm.Get(id); !ok {
		return []ID{}, true
	} else {
		return childrenSet.Keys(), true
	}
}

func (g *directedGraph[ID]) ParentsOf(id ID) ([]ID, bool) {
	if !g.checkNode(id) {
		return nil, false
	}
	var _, pm = g.guaranteeCached()
	if parentsSet, ok := pm.Get(id); !ok {
		return []ID{}, true
	} else {
		return parentsSet.Keys(), true
	}
}

func (g *directedGraph[ID]) Nodes() []ID {
	return g.baseGraph.exportNodes()
}
func (g *directedGraph[ID]) NodesLen() int {
	return g.baseGraph.countNodes()
}

func (g *directedGraph[ID]) AddNode(id ID) bool {
	atomic.AddUint64(&g.opexCount, 1)
	return g.baseGraph.addNode(id)
}

func (g *directedGraph[ID]) DeleteNode(id ID) bool {
	atomic.AddUint64(&g.opexCount, 1)
	return g.baseGraph.deleteNode(id)
}

func (g *directedGraph[ID]) CheckNode(id ID) bool {
	return g.baseGraph.checkNode(id)
}

func (g *directedGraph[ID]) GetEdge(src, tgt ID) (Weight, bool) {
	if !g.checkEdgeNodes(src, tgt) {
		return 0, false
	}
	return g.baseGraph.getEdge(src, tgt)
}

func (g *directedGraph[ID]) AddEdge(src, tgt ID, w Weight) error {
	if err := g.checkEdgeNodesWithErr(src, tgt); err != nil {
		return err
	}
	atomic.AddUint64(&g.opexCount, 1)
	if !g.baseGraph.addEdge(src, tgt, w) {
		return fmt.Errorf("%w: sourcre: %v, target: %v", ErrEdgeConflicts, src, tgt)
	}
	return nil
}

func (g *directedGraph[ID]) SetEdge(src, tgt ID, w Weight) error {
	if err := g.checkEdgeNodesWithErr(src, tgt); err != nil {
		return err
	}
	atomic.AddUint64(&g.opexCount, 1)
	g.baseGraph.setEdge(src, tgt, w)
	return nil
}

func (g *directedGraph[ID]) DeleteEdge(src, tgt ID) bool {
	if !g.checkEdgeNodes(src, tgt) {
		return false
	}
	atomic.AddUint64(&g.opexCount, 1)
	return g.baseGraph.deleteEdge(src, tgt)
}

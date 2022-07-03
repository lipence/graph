package graph

import (
	"fmt"

	"github.com/tidwall/hashmap"
)

type _DAGRuntime[ID comparable] struct {
	queue  []ID
	colors hashmap.Map[ID, Color]
}

func newTopoSortingRuntime[ID comparable](g Graph[ID]) *_DAGRuntime[ID] {
	var tsr = &_DAGRuntime[ID]{
		colors: *hashmap.New[ID, Color](g.NodesLen()),
	}
	for _, id := range g.Nodes() {
		tsr.colors.Set(id, ColorWhite)
	}
	return tsr
}

func (dr *_DAGRuntime[ID]) getColor(id ID) (Color, bool) { return dr.colors.Get(id) }
func (dr *_DAGRuntime[ID]) setColor(id ID, c Color)      { dr.colors.Set(id, c) }
func (dr *_DAGRuntime[ID]) enqueue(id ID)                { dr.queue = append(dr.queue, id) }
func (dr *_DAGRuntime[ID]) result() []ID                 { return dr.queue }

func (tr *Runtime[ID]) DAGSort(id ID) ([]ID, error) {
	var tsr = newTopoSortingRuntime(tr.graph)
	if err := tr._DAGVisit(tsr, id); err != nil {
		return nil, err
	}
	return tsr.result(), nil
}

func (tr *Runtime[ID]) DAGSortAll() ([]ID, error) {
	var tsr = newTopoSortingRuntime(tr.graph)
	for _, id := range tr.graph.Nodes() {
		if c, _ := tsr.getColor(id); c != ColorWhite {
			continue
		}
		if err := tr._DAGVisit(tsr, id); err != nil {
			return nil, err
		}
	}
	return tsr.result(), nil
}

func (tr *Runtime[ID]) _DAGVisit(tsr *_DAGRuntime[ID], id ID) error {
	switch c, _ := tsr.getColor(id); c {
	case ColorGray:
		return fmt.Errorf("not DAG")
	case ColorWhite:
		tsr.setColor(id, ColorGray)
		children, ok := tr.graph.ChildrenOf(id)
		if !ok {
			return fmt.Errorf("%w, id = %v", ErrUndefinedNode, id)
		}
		for _, w := range children {
			if err := tr._DAGVisit(tsr, w); err != nil {
				return err
			}
		}
		tsr.setColor(id, ColorBlack)
		tsr.enqueue(id)
	}
	return nil
}

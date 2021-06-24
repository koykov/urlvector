package urlvector

// Custom implementation of quick sort algorithm, special for type []vector.Node (sub-slice of vector's node array).
// Need to avoid redundant allocation when using sort.Interface.
//
// sort.Interface problem:
// <code>
// type nodes []vector.Node // type that implements sort.Interface
// ...
// children := node.Children() // get a slice of nodes to sort
// nodes := nodes(children)    // <- simple typecast, but produces an alloc (copy) due to taking address in the next line
// sort.Sort(&nodes)           // taking address
// ...
// </code>

import (
	"github.com/koykov/vector"
)

func pivot(p []vector.Node, lo, hi int) int {
	if len(p) == 0 {
		return 0
	}
	pi := &p[hi]
	i := lo - 1
	_ = p[len(p)-1]
	for j := lo; j <= hi-1; j++ {
		if p[j].KeyString() < pi.KeyString() {
			i++
			p[i].SwapWith(&p[j])
		}
	}
	if i < hi {
		p[i+1].SwapWith(&p[hi])
	}
	return i + 1
}

func quickSort(p []vector.Node, lo, hi int) {
	if lo < hi {
		pi := pivot(p, lo, hi)
		quickSort(p, lo, pi-1)
		quickSort(p, pi+1, hi)
	}
}

type nodes []vector.Node

func (n *nodes) Len() int {
	return len(*n)
}

func (n *nodes) Less(i, j int) bool {
	return (*n)[i].KeyString() < (*n)[j].KeyString()
}

func (n *nodes) Swap(i, j int) {
	(*n)[i].SwapWith(&(*n)[j])
}

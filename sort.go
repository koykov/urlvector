package urlvector

import (
	"github.com/koykov/vector"
)

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

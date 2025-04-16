package radix_tree

import "sort"

// OBJECTIVE: Represents the connection between two nodes.
// CONTENT:
//
//	label: The character at the beginning of the edge. This character determines which child node to navigate to.
//	node: The child node to which this edge is directed.
type edge struct {
	label byte
	node  *node
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e edges) Sort() {
	sort.Sort(e)
}

package router

type Node struct {
	Segment string

	Children map[string]*Node

	Parameter *Node

	Routes map[string]*Route
}

func NewNode(segment string) *Node {
	return &Node{
		Segment:  segment,
		Children: make(map[string]*Node),
		Routes:   make(map[string]*Route),
	}
}
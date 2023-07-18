package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: &node{}}
}

func (t *Tree) AddRouter(uri string, handler ControllerHandler) error {
	if !strings.HasPrefix(uri, "/") {
		return errors.New("uri must start with '/': " + uri)
	}
	if n := t.root.matchNode(uri); n != nil {
		return errors.New("router exist: " + uri)
	}

	segments := strings.Split(uri, "/")
	idx, n := 0, t.root
	for i, seg := range segments {
		if !isWildSegment(seg) {
			seg = strings.ToUpper(seg)
		}

		cnodes := n.filterNode(seg)
		var tmp *node
		for _, nn := range cnodes {
			if nn.segment == seg {
				tmp = nn
			}
		}
		if tmp == nil {
			idx = i
			break
		} else {
			n = tmp
		}
	}

	for i := idx; i < len(segments); i++ {
		nn := &node{
			segment: segments[i],
		}
		n.children = append(n.children, nn)
		n = nn
	}
	n.isLast = true
	n.handler = handler
	return nil
}

func (t *Tree) FindHandler(uri string) ControllerHandler {
	n := t.root.matchNode(strings.ToUpper(uri))
	if n != nil {
		return n.handler
	}
	return nil
}

type node struct {
	isLast   bool
	segment  string
	handler  ControllerHandler
	children []*node
}

func (n *node) filterNode(segment string) []*node {
	if len(n.children) == 0 {
		return nil
	}
	if isWildSegment(segment) {
		return n.children
	}

	nodes := make([]*node, 0)
	for _, nn := range n.children {
		if isWildSegment(nn.segment) || nn.segment == segment {
			nodes = append(nodes, nn)
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	if uri == "" || uri == "/" {
		return n
	}
	if uri[0] == '/' {
		uri = uri[1:]
	}
	segments := strings.SplitN(uri, "/", 2)

	cnodes := n.filterNode(segments[0])
	if len(cnodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		for _, nn := range cnodes {
			if nn.isLast {
				return nn
			}
		}
		return nil
	}

	for _, nn := range cnodes {
		if nextNode := nn.matchNode(segments[1]); nextNode != nil {
			return nextNode
		}
	}

	return nil
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

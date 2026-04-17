package router

import (
	"net/http"
	"strings"
)

// node is a single node in the radix tree. Each node holds a path segment,
// a map of HTTP method → handler for routes that terminate here, and child nodes
// for the next path segments.
type node struct {
	segment  string
	children []*node
	handlers map[string]http.Handler
}

// insert adds a route to the tree under the given method and pattern.
// It walks the tree segment by segment, creating nodes as needed, and stores
// the handler at the leaf node for the pattern.
func (n *node) insert(method, pattern string, handler http.Handler) {
	if pattern != "" {
		patternSegments := strings.SplitSeq(pattern, "/")
		for seg := range patternSegments {
			var next *node
			for _, child := range n.children {
				if child.segment == seg {
					next = child
					break
				}
			}
			if next == nil {
				next = &node{segment: seg}
				n.children = append(n.children, next)
			}
			n = next

		}
		if n.handlers == nil {
			n.handlers = make(map[string]http.Handler)
		}
		n.handlers[method] = handler
	}
}

// search walks the tree to find the node matching the given path.
// Returns the matched node and a map of captured URL parameters, or nil if no match.
// Exact segments take priority over named params, which take priority over wildcards.
func (n *node) search(path string) (*node, map[string]string) {
	if path == "" {
		return nil, nil
	}
	segments := strings.Split(path, "/")
	params := make(map[string]string)
	for i, seg := range segments {
		var matched *node
		for _, child := range n.children {
			if child.segment == seg {
				matched = child
				break
			} else if strings.HasPrefix(child.segment, "{") && strings.HasSuffix(child.segment, "}") && !strings.HasSuffix(child.segment, ":*}") {
				name := child.segment[1 : len(child.segment)-1]
				params[name] = seg
				matched = child
				break
			} else if child.segment == "*" {
				params["*"] = strings.Join(segments[i:], "/")
				return child, params
			} else if strings.HasSuffix(child.segment, ":*}") {
				name := child.segment[1 : len(child.segment)-3]
				params[name] = strings.Join(segments[i:], "/")
				return child, params
			}
		}
		if matched == nil {
			return nil, nil
		}
		n = matched
	}
	return n, params
}

package traversal

import (
	"fmt"
	"tubes2/backend/internal/parser"
	"tubes2/backend/internal/selector"
)

func DFS(limit int, root *parser.Node, selectors []string) ([]*parser.Node, []*parser.Node, []string) {
	var matches []*parser.Node
	var visited []*parser.Node
	var log []string
	dfsHelper(root, selectors, limit, &matches, &visited, &log)
	return matches, visited, log
}


func dfsHelper(node *parser.Node, selectors []string, limit int, matches *[]*parser.Node, visited *[]*parser.Node, log *[]string) bool {
	if node == nil {
		return false
	}

	*visited = append(*visited, node)
	*log = append(*log, fmt.Sprintf("Visit: %s (depth=%d)", nodeLabel(node), node.Depth))

	if selector.ComboCombinator(node, selectors) {
		*matches = append(*matches, node)
		*log = append(*log, fmt.Sprintf("Match: %s (depth=%d)", nodeLabel(node), node.Depth))
		if limit > 0 && len(*matches) == limit {
			return true
		}
	}

	for _, child := range node.Children {
		if dfsHelper(child, selectors, limit, matches, visited, log) {
			return true
		}
	}
	return false
}
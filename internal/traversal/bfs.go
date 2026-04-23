package traversal

import (
	"fmt"
	"tubes2/backend/internal/parser"
	"tubes2/backend/internal/selector"
)

func BFS(limit int, root *parser.Node, selectors []string) ([]*parser.Node, []*parser.Node, []string) {
	if root == nil {
		return nil, nil, nil
	}

	queue := []*parser.Node{root}
	var visited []*parser.Node
	var matches []*parser.Node
	var log []string
	count := 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		visited = append(visited, node)
		log = append(log, fmt.Sprintf("Visit: %s (depth=%d)", nodeLabel(node), node.Depth))

		if selector.ComboCombinator(node, selectors) {
			matches = append(matches, node)
			count++
			log = append(log, fmt.Sprintf("Match: %s (depth=%d)", nodeLabel(node), node.Depth))
			if limit > 0 && count == limit {
				return matches, visited, log
			}
		}

		for _, child := range node.Children {
			queue = append(queue, child)
		}
	}
	return matches, visited, log
}
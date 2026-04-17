package traversal

import (
	"tubes2/backend/internal/parser"
	"tubes2/backend/internal/selector"
)

func BFS(limit int, root *parser.Node, selectors []string) ([]parser.Node, []parser.Node) {
	if root == nil {
		return nil, nil
	}
	queue := []*parser.Node{root}
	var visited []parser.Node
	var res []parser.Node
	for count := 0; len(queue) > 0; {
		visited = append(visited, *queue[0])
		if selector.ComboCombinator(root, selectors) {
			res = append(res, *queue[0])
			count++
			if count == limit {
				return res, visited
			}
		}
		queue = queue[1:]
		for _, child := range queue[0].Children {
			queue = append(queue, child)
		}
	}
	return res, visited
}

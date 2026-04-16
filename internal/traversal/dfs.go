package traversal

import "tubes2/backend/internal/parser"

func DFS(count, limit int, root *parser.Node) ([]parser.Node, []parser.Node) { // call with count = 0
	var visited []parser.Node
	var res []parser.Node
	visited = append(visited, *root)
	if true /*TODO: selectors*/ {
		if count == limit {
			return res, visited
		}
		res = append(res, *root)
		count++
	}
	for _, child := range root.Children {
		if count == limit {
			break
		}
		res1, visited1 := DFS(count, limit, child)
		res, visited = append(res, res1...), append(visited, visited1...)
	}
	return res, visited
}

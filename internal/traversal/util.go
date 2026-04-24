package traversal

import "tubes2/backend/internal/parser"

func nodeLabel(n *parser.Node) string {
	label := "<" + n.Tag
	if n.ID != "" {
		label += "#" + n.ID
	}
	for _, cls := range n.Classes {
		label += "." + cls
	}
	label += ">"
	return label
}
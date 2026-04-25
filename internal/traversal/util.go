package traversal

import "tubes2/backend/internal/parser"

func nodeLabel(n *parser.Node) string {
	if n.Tag == "#text" {
		t := n.Text
		if len(t) > 15 {
			t = t[:10] + "..."
		}
		return "<#text: \"" + t + "\">"
	}
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
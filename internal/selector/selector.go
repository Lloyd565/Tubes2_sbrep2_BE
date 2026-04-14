package selector

import "tubes2/backend/internal/parser"

// Combinators
func ChildCombinator(a, b *parser.Node) bool { // Symbol: ">"
	return a == b.Parent
}

func DescendentCombinator(a, b *parser.Node) bool { // Symbol: " "
	depth := b.Depth
	t := b
	for depth > a.Depth {
		t = t.Parent
		depth--
	}
	return a == t
}

func AdjSiblingCombinator(a, b *parser.Node) bool { // Symbol: "+"
	parent := a.Parent
	if parent != b.Parent {
		return false
	}
	idx := 0
	for parent.Children[idx] != a {
		idx++
	}
	return parent.Children[idx+1] == b
}

func GenSiblingCombinator(a, b *parser.Node) bool { // Symbol: "~"
	parent := a.Parent
	if parent != b.Parent {
		return false
	}
	idx := 0
	for parent.Children[idx] != a {
		idx++
	}
	for i := idx; i < len(parent.Children); i++ {
		if parent.Children[i] == b {
			return true
		}
	}
	return false
}

// Selectors
func TagSelector(a *parser.Node, tag string) bool {
	return a.Tag == tag
}

func ClassSelector(a *parser.Node, classes []string) bool {
	if len(classes) == 0 {
		return true
	}
	if len(a.Classes) < len(classes) {
		return false
	}
	temp := classes
	for _, cls1 := range a.Classes {
		for i, cls2 := range temp {
			if cls1 == cls2 {
				if i == len(temp)-1 {
					temp = temp[:i]
					break
				}
				temp = append(temp[:i-1], temp[i+1:]...)
				break
			}
		}
		if len(temp) == 0 {
			return true
		}
	}
	return false
}

func AttrSelector(a *parser.Node, attr string, val string)

func IDSelector(a *parser.Node, id string) bool {
	return a.ID == id
}

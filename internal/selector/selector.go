package selector

import "tubes2/backend/internal/parser"

// Combinators
func ChildCombinator(a, b *parser.Node) bool {
	return a == b.Parent
}

func DescendentCombinator(a, b *parser.Node) bool {
	current := b.Parent
	for current != nil {
		if current == a {
			return true
		}
		current = current.Parent
	}
	return false
	// depth := b.Depth
	// t := b
	// for depth > a.Depth {
	// 	t = t.Parent
	// 	depth--
	// }
	// return a == t
}

func AdjSiblingCombinator(a, b *parser.Node) bool {
	if a.Parent == nil || a.Parent != b.Parent {
		return false
	}
	children := a.Parent.Children
	for i, child := range children {
		if child == a {
			if i+1 >= len(children) {
				return false // a adalah anak terakhir
			}
			return children[i+1] == b
		}
	}
	return false
}

func GenSiblingCombinator(a, b *parser.Node) bool {
	parent := a.Parent
	if parent != b.Parent {
		return false
	}
	idx := 0
	for parent.Children[idx] != a {
		idx++
	}
	for i := idx + 1; i < len(parent.Children); i++ {
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
	for _, target := range classes {
		found := false
		for _, has := range a.Classes {
			if has == target {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func AttrSelector(a *parser.Node, attr string, val string) bool {
	v, exist := a.Attr[attr]
	if !exist {
		return false
	}
	if val == "" {
		return true
	}
	return v == val
}

func IDSelector(a *parser.Node, id string) bool {
	return a.ID == id
}

// Matching Selector
func MatchSelector(node *parser.Node)

package selector

import (
	"tubes2/backend/internal/parser"
)

func TagSelector(a *parser.Node, tag string) bool {
	return a.Tag == tag || tag == "*"
}

func ClassSelector(a *parser.Node, classes []string) bool {
	for _, cls := range classes {
		found := false
		for _, nodeClass := range a.Classes {
			if cls == nodeClass {
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

func AttrSelector(a *parser.Node, attrs map[string]string) bool {
	for attr, val := range attrs {
		nodeVal, exists := a.Attributes[attr]
		if !exists {
			return false
		}
		if val != "" && val != nodeVal {
			return false
		}
	}
	return true
}

func IDSelector(a *parser.Node, id string) bool {
	return a.ID == id
}

func ComboSelector(a *parser.Node, selectors []string) bool {
	if a == nil {
		return false
	}
	if len(selectors) == 0 {
		return true
	}
	var classes []string
	attrs := make(map[string]string)
	for _, str := range selectors {
		if str == "*" {
			return true
		}
		switch str[0] {
		case '.':
			classes = append(classes, str[1:])
		case '#':
			if !IDSelector(a, str[1:]) {
				return false
			}
		case '[':
			attr, val := parser.AttrParser(str)
			attrs[attr] = val
		default:
			if !TagSelector(a, str) {
				return false
			}
		}
	}
	return ClassSelector(a, classes) && AttrSelector(a, attrs)
}

func ComboCombinator(a *parser.Node, selector []string) bool {
	start := 0
	for i, s := range selector {
		if s == "," {
			if matchSingle(a, selector[start:i]) {
				return true
			}
			start = i + 1
		}
	}
	return matchSingle(a, selector[start:])
}

func matchSingle(a *parser.Node, selector []string) bool {
	if len(selector) == 0 {
		return false
	}

	i := len(selector) - 1
	if !ComboSelector(a, parser.SelectorParser(selector[i])) {
		return false
	}

	current := a
	i--

	for i >= 0 {
		tok := selector[i]
		switch tok[0] {
		case '>':
			i--
			if i < 0 {
				return false
			}
			if current.Parent == nil || !ComboSelector(current.Parent, parser.SelectorParser(selector[i])) {
				return false
			}
			current = current.Parent

		case '+':
			i--
			if i < 0 {
				return false
			}
			if current.Parent == nil {
				return false
			}
			pos := -1
			for j, node := range current.Parent.Children {
				if node == current {
					pos = j
					break
				}
			}
			if pos <= 0 || !ComboSelector(current.Parent.Children[pos-1], parser.SelectorParser(selector[i])) {
				return false
			}
			current = current.Parent.Children[pos-1]

		case '~':
			i--
			if i < 0 {
				return false
			}
			if current.Parent == nil {
				return false
			}
			found := false
			for _, node := range current.Parent.Children {
				if node == current {
					break
				}
				if ComboSelector(node, parser.SelectorParser(selector[i])) {
					found = true
					current = node
				}
			}
			if !found {
				return false
			}

		default:
			temp := current.Parent
			for temp != nil {
				if ComboSelector(temp, parser.SelectorParser(tok)) {
					current = temp
					break
				}
				temp = temp.Parent
			}
			if temp == nil {
				return false
			}
		}
		i--
	}
	return true
}
package selector

import (
	"tubes2/backend/internal/parser"
)

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

// combo

func ComboSelector(a *parser.Node, selectors []string) bool {
	var classes []string
	for _, str := range selectors {
		switch str[0] {
		case '.':
			classes = append(classes, str[1:len(str)-1])
		case '#':
			if IDSelector(a, str[1:len(str)-1]) {
				return false
			}
		case '[': // TODO attribute selector
		default:
			if TagSelector(a, str) {
				return false
			}
		}
	}
	return ClassSelector(a, classes)
}

func ComboCombinator(a *parser.Node, selector []string) bool {
	i := len(selector) - 1
	selectors := parser.SelectorParser(selector[i])
	if !ComboSelector(a, selectors) {
		return false
	}
	for i := i - 1; i >= 0; i-- {
		switch selector[i][0] {
		case '>':
			if i == 0 {
				return false
			}
			selectors = parser.SelectorParser(selector[i-1])
			if ComboSelector(a.Parent, selectors) {
				return false
			}
		case '+':
			if i == 0 {
				return false
			}
			selectors = parser.SelectorParser(selector[i-1])
			for i, node := range a.Parent.Children {
				if node == a {
					if i == 0 {
						return false
					}
					if !ComboSelector(a.Parent.Children[i-1], selectors) {
						return false
					}
					break
				}
			}
		case '~':
			if i == 0 {
				return false
			}
			selectors = parser.SelectorParser(selector[i-1])
			for _, node := range a.Parent.Children {
				if node == a {
					return false
				}
				if ComboSelector(node, selectors) {
					break
				}
			}
		default:
			selectors = parser.SelectorParser(selector[i])
			temp := a.Parent
			for ComboSelector(temp, selectors) {
				if temp == nil {
					return false
				}
				temp = temp.Parent
			}
		}
	}
	return true
}

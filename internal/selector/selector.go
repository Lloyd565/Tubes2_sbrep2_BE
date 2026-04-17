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
	for i, cls1 := range classes {
		found := false
		for _, cls2 := range a.Classes {
			if cls1 == cls2 {
				found = true
				if i == len(classes)-1 {
					classes = classes[:i]
					break
				}
				classes = append(classes[:i-1], classes[i+1:]...)
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
		found := false
		for attr1, val1 := range a.Attributes {
			if attr == attr1 {
				if val == "" || val == val1 {
					delete(attrs, attr)
					found = true
					break
				} else {
					return false
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func IDSelector(a *parser.Node, id string) bool {
	return a.ID == id
}

// combo

func ComboSelector(a *parser.Node, selectors []string) bool {
	var classes []string
	var attrs map[string]string
	for _, str := range selectors {
		switch str[0] {
		case '.':
			classes = append(classes, str[1:len(str)-1])
		case '#':
			if IDSelector(a, str[1:len(str)-1]) {
				return false
			}
		case '[':
			attr, val := parser.AttrParser(str)
			attrs[attr] = val
		default:
			if TagSelector(a, str) {
				return false
			}
		}
	}
	return ClassSelector(a, classes) && AttrSelector(a, attrs)
}

func ComboCombinator(a *parser.Node, selector []string) bool {
	i := len(selector) - 1
	selectors := parser.SelectorParser(selector[i])
	if !ComboSelector(a, selectors) {
		return false
	}
	if selector[0][0] == '>' || selector[0][0] == '+' || selector[0][0] == '~' || selector[0][0] == ',' {
		return false
	}
	valid := true
	for i := i - 1; i >= 0; i-- {
		switch selector[i][0] {
		case '>':
			selectors = parser.SelectorParser(selector[i-1])
			if ComboSelector(a.Parent, selectors) {
				valid = false
			}
		case '+':
			selectors = parser.SelectorParser(selector[i-1])
			for i, node := range a.Parent.Children {
				if node == a {
					if i == 0 {
						valid = false
					}
					if !ComboSelector(a.Parent.Children[i-1], selectors) {
						valid = false
					}
					break
				}
			}
		case '~':
			selectors = parser.SelectorParser(selector[i-1])
			for _, node := range a.Parent.Children {
				if node == a {
					valid = false
				}
				if ComboSelector(node, selectors) {
					break
				}
			}
		case ',':
			if valid {
				return true
			}
			valid = true
		default:
			selectors = parser.SelectorParser(selector[i])
			temp := a.Parent
			for ComboSelector(temp, selectors) {
				if temp == nil {
					valid = false
				}
				temp = temp.Parent
			}
		}
	}
	return false
}

package parser

import (
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	Tag        string
	ID         string
	Classes    []string
	Parent     *Node
	Depth      int
	Children   []*Node
	Attributes map[string]string
}

func Parse(htmlStr string) (*Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var root *html.Node
	var findHTML func(*html.Node)
	findHTML = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "html" {
			root = n
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			findHTML(child)
		}
	}
	findHTML(doc)

	return buildNode(root, nil, 0), nil
}

func buildNode(n *html.Node, parent *Node, depth int) *Node {
	if n == nil {
		return nil
	}

	node := &Node{
		Tag:        n.Data,
		Parent:     parent,
		Depth:      depth,
		Attributes: make(map[string]string),
	}

	for _, attr := range n.Attr {
		node.Attributes[attr.Key] = attr.Val
		if attr.Key == "id" {
			node.ID = attr.Val
		}
		if attr.Key == "class" {
			node.Classes = strings.Fields(attr.Val)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			child := buildNode(c, node, depth+1)
			node.Children = append(node.Children, child)
		}
	}

	return node
}

func CombinatorParser(selector string) []string {
	var res []string
	idx := 0
	for i, str := range selector {
		if str == ' ' {
			if selector[i-1] == ',' {
				res = append(res, selector[idx:i-2])
				res = append(res, selector[i-1:i-1])
			} else {
				res = append(res, selector[idx:i-1])
			}
			idx = i + 1
		}
	}
	return res
}

func SelectorParser(selector string) []string {
	var res []string
	idx := 0
	for i, str := range selector {
		if (str == '.' || str == '[' || str == '#') && i != 0 {
			res = append(res, selector[idx:i-1])
			idx = i
		}
	}
	return res
}

func AttrParser(attrs string) (string, string) {
	attrs = attrs[1 : len(attrs)-2]
	idx := 0
	var attr string
	for i, str := range attrs {
		if str == '=' {
			attr = attrs[idx : i-1]
			idx = i + 1
			break
		}
		if i == len(attrs)-1 {
			return attrs, ""
		}
	}
	return attr, attrs[idx:]
}

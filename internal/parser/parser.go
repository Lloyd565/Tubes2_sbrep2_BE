package parser

import (
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	Tag      string
	ID       string
	Classes  []string
	Parent   *Node
	Depth    int
	Children []*Node
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
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findHTML(c)
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
		Tag:    n.Data,
		Parent: parent,
		Depth:  depth,
	}
	for _, attr := range n.Attr {
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
			res = append(res, selector[idx:i])
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

// TODO: attribute parser

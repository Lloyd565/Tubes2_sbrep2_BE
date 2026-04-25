package parser

import (
	"strings"

	"golang.org/x/net/html"
)

var selfClosing = map[string]bool{
	"area": true, "base": true, "br": true, "col": true,
	"embed": true, "hr": true, "img": true, "input": true,
	"link": true, "meta": true, "param": true, "source": true,
	"track": true, "wbr": true,
}

type Node struct {
	Tag        string
	ID         string
	Classes    []string
	Parent     *Node
	Depth      int
	NodeNum int
	Children   []*Node
	Attributes map[string]string
}

func Parse(htmlStr string) (*Node, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlStr))

	virtualRoot := &Node{Tag: "#document", Depth: -1, Attributes: make(map[string]string)}
	stack := []*Node{virtualRoot}
	nodeCnt := 1

	for {
		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			break
		}

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			rawName, hasAttr := tokenizer.TagName()
			tagName := string(rawName)

			parent := stack[len(stack)-1]
			node := &Node{
				Tag:        tagName,
				Parent:     parent,
				Depth:      parent.Depth + 1,
				NodeNum: nodeCnt,
				Attributes: make(map[string]string),
			}
			nodeCnt++

			for hasAttr {
				var k, v []byte
				k, v, hasAttr = tokenizer.TagAttr()
				key, val := string(k), string(v)
				node.Attributes[key] = val
				switch key {
				case "id":
					node.ID = val
				case "class":
					node.Classes = strings.Fields(val)
				}
			}

			parent.Children = append(parent.Children, node)

			isSelfClose := tt == html.SelfClosingTagToken || selfClosing[tagName]
			if !isSelfClose {
				stack = append(stack, node)
			}

		case html.EndTagToken:
			rawName, _ := tokenizer.TagName()
			tagName := string(rawName)

			for len(stack) > 1 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.Tag == tagName {
					break
				}
			}
		}
	}

	for _, child := range virtualRoot.Children {
		if child.Tag == "html" {
			child.Parent = nil
			return child, nil
		}
	}
	if len(virtualRoot.Children) > 0 {
		root := virtualRoot.Children[0]
		root.Parent = nil
		return root, nil
	}

	return virtualRoot, nil
}

func CombinatorParser(selector string) []string {
	var res []string
	var cur strings.Builder

	flushCur := func() {
		tok := strings.TrimSpace(cur.String())
		if tok != "" {
			res = append(res, tok)
		}
		cur.Reset()
	}

	inBrt := false
	i := 0
	for i < len(selector) {
		ch := selector[i]

		if ch == '[' {
			inBrt = true
		} else if ch == ']' {
			inBrt = false
		}

		switch {
		case inBrt:
			cur.WriteByte(ch)

		case ch == ' ':
			flushCur()
			for i+1 < len(selector) && selector[i+1] == ' ' {
				i++
			}

		case ch == '>' || ch == '+' || ch == '~':
			flushCur()
			res = append(res, string(ch))

		case ch == ',':
			flushCur()
			res = append(res, ",")

		default:
			cur.WriteByte(ch)
		}
		i++
	}
	flushCur()
	return res
}

func SelectorParser(selector string) []string {
	var res []string
	idx := 0
	for i, ch := range selector {
		if (ch == '.' || ch == '[' || ch == '#') && i != 0 {
			res = append(res, selector[idx:i])
			idx = i
		}
	}
	res = append(res, selector[idx:])
	return res
}

func AttrParser(attrs string) (string, string) {
	attrs = attrs[1 : len(attrs)-1]
	for i, ch := range attrs {
		if ch == '=' {
			val := attrs[i+1:]
			if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') && val[len(val)-1] == val[0] {
				val = val[1 : len(val)-1]
			}
			return attrs[:i], val
		}
	}
	return attrs, ""
}
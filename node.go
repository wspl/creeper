package creeper

import (
	"regexp"
	"strings"
)

var (
	rx_node = regexp.MustCompile(`^(\s*)([a-zA-Z0-9-_]+)\s*(\[]|\*)?\s*:\s*(.+)$`)
	rx_page = regexp.MustCompile(`^[a-zA-Z-]`)
	rx_fun  = regexp.MustCompile(`^($|.)`)
)

type Node struct {
	Raw     string
	Creeper *Creeper

	Name      string
	IsArray   bool
	IsPrimary bool
	IndentLen int

	Page *Page
	Fun  *Fun

	Index int

	Sn map[int]string

	PrevNode       *Node
	NextNode       *Node
	ParentNode     *Node
	FirstChildNode *Node
	LastChildNode  *Node
}

func (n *Node) Inc() {
	n.Index++
}

func (n *Node) Reset() {
	n.Index = 0
}

func (n *Node) Filter(cb func(*Node) bool) []*Node {
	nodes := []*Node{}
	for node := n; node != nil; node = node.PrevNode {
		if cb(node) {
			nodes = append(nodes, node)
		}
	}
	if n.NextNode != nil {
		for node := n.NextNode; node != nil; node = node.NextNode {
			if cb(node) {
				nodes = append(nodes, node)
			}
		}
	}
	return nodes
}

func (n *Node) SearchRef(name string) *Node {
	for node := n.PrevNode; node != nil; node = node.PrevNode {
		if node.Name == name {
			return node
		}
	}
	for node := n.ParentNode; node != nil; node = node.ParentNode {
		for node := node; node != nil; node = node.PrevNode {
			if node.Name == name {
				return node
			}
		}
	}
	return nil
}

func (n *Node) SearchFlatScope(name string) *Node {
	ns := n.Filter(func(n *Node) bool {
		return n.Name == name
	})
	if len(ns) > 0 {
		return ns[0]
	} else {
		return nil
	}
}

func (n *Node) Search(name string) *Node {
	if n.FirstChildNode == nil {
		return nil
	}
	ns := n.FirstChildNode.Filter(func(n *Node) bool {
		return n.Name == name
	})
	if len(ns) > 0 {
		return ns[0]
	} else {
		return nil
	}
}

func (n *Node) Primary() *Node {
	if n.FirstChildNode == nil {
		return nil
	}
	ns := n.FirstChildNode.Filter(func(n *Node) bool {
		return n.IsPrimary
	})
	if len(ns) > 0 {
		return ns[0]
	} else {
		return nil
	}
}

func (n *Node) Value() (string, error) {
	if n.Fun != nil {
		return n.Fun.Invoke()
	} else {
		return n.Page.Body()
	}
}

func ParseNode(ln []string) []*Node {
	nodes := []*Node{}
	var justNode *Node
	for i, l := range ln {
		node := new(Node)
		node.Raw = l
		node.Sn = splitNode(l)
		node.Name = node.Sn[1]
		node.IndentLen = len(node.Sn[0])

		switch node.Sn[4] {
		case "[]":
			node.IsArray = true
		case "*":
			node.IsPrimary = true
		}

		if justNode != nil {
			if node.IndentLen == justNode.IndentLen {
				node.ParentNode = justNode.ParentNode
				node.PrevNode = justNode
				justNode.NextNode = node
			} else if node.IndentLen > justNode.IndentLen {
				node.ParentNode = justNode
				justNode.FirstChildNode = node
			} else if node.IndentLen < justNode.IndentLen {
				justNode.ParentNode.NextNode = node
				justNode.ParentNode.LastChildNode = justNode
				node.PrevNode = justNode.ParentNode
				node.ParentNode = justNode.ParentNode.ParentNode
			}
			if i == len(ln)-1 {
				n := node
				for n.ParentNode != nil {
					n.ParentNode.LastChildNode = n
					n = n.ParentNode
				}
			}
		}

		nodes = append(nodes, node)
		justNode = node
	}
	for _, node := range nodes {
		if len(node.Sn[2]) > 0 {
			node.Page = ParsePage(node, node.Sn[2])
		}
		if len(node.Sn[3]) > 0 {
			if node.Sn[3][0] == '.' {
				node.Fun = ParseFun(node, "$" + node.Sn[3])
			} else {
				node.Fun = ParseFun(node, node.Sn[3])
			}
		}
	}
	return nodes
}

func splitNode(s string) map[int]string {
	rs := map[int]string{}

	sArr := rx_node.FindAllStringSubmatch(s, -1)
	// indent
	rs[0] = sArr[0][1]
	// name
	rs[1] = strings.TrimSpace(sArr[0][2])
	// is array or primary key
	rs[4] = sArr[0][3]

	pfArr := strings.Split(sArr[0][4], "->")
	if len(pfArr) == 1 {
		if rx_page.MatchString(pfArr[0]) {
			rs[2] = strings.TrimSpace(pfArr[0])
		} else if rx_fun.MatchString(pfArr[0]) {
			rs[3] = strings.TrimSpace(pfArr[0])
		}
	} else {
		rs[2] = strings.TrimSpace(pfArr[0])
		rs[3] = strings.TrimSpace(pfArr[1])
	}

	return rs
}

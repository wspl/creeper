package creeper

import (
	"regexp"
	"strings"
)

var (
	rx_node = regexp.MustCompile(`^(\s*)([a-zA-Z0-9-_]+)\s*(\[])?\s*:\s*(.+)$`)
	rx_page = regexp.MustCompile(`^[a-zA-Z-]`)
	rx_fun  = regexp.MustCompile(`^($|.)`)
)

type Node struct {
	Raw string
	Creeper *Creeper

	Name      string
	IsArray   bool
	IndentLen int

	Page *Page
	Fun  *Fun

	Index int

	PrevNode   *Node
	NextNode   *Node
	ParentNode *Node
	FirstChildNode *Node
	LastChildNode *Node
}

func (n *Node) Inc() {
	n.Index++
}

func (n *Node) Reset() {
	n.Index = 0
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
	for node := n; node != nil; node = node.PrevNode {
		if node.Name == name {
			return node
		}
	}
	for node := n.NextNode; node != nil; node = node.NextNode {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (n *Node) Value() string {
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
		sn, isArray := splitNode(l)
		node.Name = sn[1]
		node.IsArray = isArray
		node.IndentLen = len(sn[0])

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
			if i == len(ln) - 1 {
				n := node
				for n.ParentNode != nil {
					n.ParentNode.LastChildNode = n
					n = n.ParentNode
				}
			}
		}

		if len(sn[2]) > 0 {
			node.Page = ParsePage(node, sn[2])
		}
		if len(sn[3]) > 0 {
			if sn[3][0] == '.' {
				node.Fun = ParseFun(node, node.ParentNode.Fun.Raw + sn[3])
				node.Page = node.ParentNode.Page
			} else {
				node.Fun = ParseFun(node, sn[3])
			}
		}

		nodes = append(nodes, node)
		justNode = node
	}
	return nodes
}

func splitNode(s string) (map[int]string, bool) {
	rs := map[int]string{}
	isArray := false

	sArr := rx_node.FindAllStringSubmatch(s, -1)
	// indent
	rs[0] = sArr[0][1]
	// name
	rs[1] = strings.TrimSpace(sArr[0][2])
	// is array
	if len(sArr[0][3]) == 2 {
		isArray = true
	}

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

	return rs, isArray
}
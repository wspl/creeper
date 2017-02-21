package creeper

import (
	"regexp"
	"strings"
)

type Formatted struct {
	Raw string

	Towns []*Town
	Nodes []*Node
}

var (
	rx_isTown = regexp.MustCompile(`^\s*[a-zA-Z][a-zA-Z_-]{0,31}\s*(\(|=)`)
	rx_isNode = regexp.MustCompile(`^\s*[a-zA-Z_-]{1,32}(\[])?:`)
)

func Formatting(s string) *Formatted {
	ln := strings.Split(s, "\n")

	ln = stringsFilter(ln, func(s string) bool {
		l := strings.TrimSpace(s)
		return len(l) > 0 && l[0] != '#'
	})
	ln = stringsMap(ln, func(s string) string {
		return strings.Replace(s, "\r", "", -1)
	})
	ln = linesCombine(ln)

	townLn := stringsMatch(ln, rx_isTown)
	nodeLn := stringsMatch(ln, rx_isNode)

	towns := ParseTown(townLn)
	nodes := ParseNode(nodeLn)

	return &Formatted{
		Raw:   s,
		Towns: towns,
		Nodes: nodes,
	}
}

func linesCombine(l []string) []string {
	nl := []string{}
	for i, s := range l {
		if strings.TrimSpace(s)[0] != '.' {
			s := s
			for i := i + 1; i < len(l); i++ {
				ns := strings.TrimSpace(l[i])
				if ns[0] != '.' {
					break
				}
				s = s + ns
			}
			nl = append(nl, s)
		}
	}
	return nl
}

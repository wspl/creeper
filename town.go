package creeper

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

var (
	rx_townName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]{0,31}`)
)

func Town_New() *Town {
	t := new(Town)
	t.Params = map[string]string{}
	return t
}

type Town struct {
	Raw     string
	Creeper *Creeper
	Node    *Node

	Name     string
	Params   map[string]string
	Template string
}

func (t *Town) Value() string {
	s := t.Template
	for k, v := range t.Params {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}

func (t *Town) HasParam(k string) bool {
	_, e := t.Params[k]
	return e
}

func (t *Town) Set(k string, v string) bool {
	e := t.HasParam(k)
	t.Params[k] = v
	return e
}

func (t *Town) Get(k string) (string, bool) {
	v, e := t.Params[k]
	return v, e
}

func (t *Town) PreSet(k string) bool {
	return t.Set(k, "")
}

func (t *Town) Attach() bool {
	var rt *Town
	for _, tt := range t.Node.Creeper.Towns {
		if t.Name == tt.Name {
			rt = tt
		}
	}
	if rt == nil {
		return false
	}
	t.Template = rt.Template
	tp := t.Params
	t.Params = rt.Params
	for k, v := range tp {
		t.Params[k] = v
	}
	return true
}

func ParseTown(ln []string) []*Town {
	towns := []*Town{}
	for _, l := range ln {
		towns = append(towns, ParseTownLine(l))
	}
	return towns
}

func ParseTownLine(l string) *Town {
	town := Town_New()
	town.Raw = l
	sr := rx_townName.FindAllString(l, -1)
	town.Name = sr[0]
	ls := strings.TrimSpace(l[len(sr[0]):])
	if len(ls) > 0 {
		if ls[0] == '(' {
			pa, ep := parseParams(ls)
			for k, v := range pa {
				town.Set(k, v)
			}
			ls = ls[ep:]
		}
		if len(ls) > 1 {
			td := strings.TrimSpace(strings.TrimSpace(ls)[1:])
			if len(td) > 0 && td[0] == '=' {
				td = strings.TrimSpace(td[1:])
			}
			town.Template = trimTownValue(td)
		}
	}
	return town
}

func trimTownValue(s string) string {
	s = strings.TrimSpace(s)
	if s[0] == '"' || s[0] == '`' {
		s = s[1 : len(s)-1]
		if s[0] == '"' {
			s = strings.Replace(s, `\\`, `\`, -1)
		}
	}
	return s
}

// start with "(", will return params map and end pos.
// all params string type:
// (key1 = 0, key2 = "str_exam\"ple", key3 = `exp_\`example\n`)
// (key1 = 0, key2, key3)
// (key1, key2, key3)
// ("str_exam\"ple", /exp_\/example\n/, 2)
func parseParams(s string) (map[string]string, int) {
	endPos := -1

	kvMap := map[string]string{}
	pK := ""
	pIsK := false

	var sb bytes.Buffer

	inKey := false
	inStr := false // "example"
	inExp := false // `example`
	inStd := false //  example

	for i, c := range s {
		cso := func(o int) int32 {
			oi := i + o
			if oi >= 0 && oi < len(s) {
				return rune(s[oi])
			}
			return 0
		}
		co := func(o int) int32 {
			if i+o < 0 || i+0 >= len(s) {
				return 0
			}
			if o < 0 {
				j := i
				for j >= 0 && o != 0 {
					j--
					if !unicode.IsSpace(rune(s[j])) {
						o++
					}

				}
				return rune(s[j])
			} else if o > 0 {
				j := i
				for j < len(s)-1 && o != 0 {
					j++
					if !unicode.IsSpace(rune(s[j])) {
						o--
					}
				}
				return rune(s[j])
			} else {
				return rune(s[i])
			}
			return 0
		}

		if i == 0 && c != '(' {
			return nil, -1
		}

		if !inExp && !inStr && !inStd {
			if (co(-1) == '(' || co(-1) == ',') && (unicode.IsLetter(c) || c == '@') {
				inKey = true
			} else if (co(-1) == '=' || co(-1) == ',' || co(-1) == '(') &&
				!unicode.IsSpace(c) && c != '"' && c != '`' {
				inStd = true
			} else if co(-2) == '=' || co(-2) == ',' || co(-2) == '(' {
				switch co(-1) {
				case '"':
					inStr = true
				case '`':
					inExp = true
				}
			}
		}

		if inKey || inExp || inStd || inStr {
			sb.WriteRune(c)
		}

		if !inExp && !inStd && !inStr && !inKey && c == ')' {
			endPos = i
		}

		if c != '\\' {
			if inKey && (co(1) == ',' || co(1) == ')' || co(1) == '=') {
				inKey = false
				//println("key: ", sb.String())
				pK = strings.TrimSpace(sb.String())
				kvMap[pK] = ""
				if co(1) != ',' {
					pIsK = true
				}
				sb.Reset()
			} else if inStr && cso(1) == '"' {
				inStr = false
				//println("str: ", sb.String())
				s := strings.TrimSpace(sb.String())
				s = strings.Replace(s, `\\`, `\`, -1)
				if pIsK {
					kvMap[pK] = strings.TrimSpace(sb.String())
				} else {
					kvMap[s] = ""
				}
				pIsK = false
				sb.Reset()
			} else if inExp && cso(1) == '`' {
				inExp = false
				//println("exp: ", sb.String())
				s := strings.TrimSpace(sb.String())
				if pIsK {
					kvMap[pK] = strings.TrimSpace(sb.String())
				} else {
					kvMap[s] = ""
				}
				pIsK = false
				sb.Reset()
			} else if inStd && (co(1) == ',' || co(1) == ')') {
				inStd = false
				//println("std: ", sb.String())
				s := strings.TrimSpace(sb.String())
				if pIsK {
					kvMap[pK] = strings.TrimSpace(sb.String())
				} else {
					kvMap[s] = ""
				}
				pIsK = false
				sb.Reset()
			}
		}

		if endPos > -1 {
			break
		}
	}
	return kvMap, endPos
}

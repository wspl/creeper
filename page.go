package creeper

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type Page struct {
	Raw  string
	Node *Node

	Town     *Town
	Ref      string

	NextMode bool
	NextUrl string
	NextPendingUrl string
	NextReady bool
	NextNoMore bool

	Index int
}

func (p *Page) Inc() {
	p.Index++
}

func (p *Page) Url() (string, error) {
	if p.NextMode {
		return p.NextUrl, nil
	}
	if p.Town != nil {
		p.Town.Attach()
		if p.Index > -1 {
			p.Town.Set("@page", strconv.Itoa(p.Index))
		} else {
			i, e := p.Town.Get("@page")
			if e {
				i64, _ := strconv.ParseInt(i, 10, 32)
				p.Index = int(i64)
			}
		}
		for k, v := range p.Town.Params {
			if len(v) > 0 && v[0] == '_' {
				n := p.Node.SearchRef(v)
				v, err := n.Value()
				if err != nil {
					return "", err
				}
				p.Town.Set(k, v)
			}
		}
		return p.Town.Value(), nil
	} else {
		return p.Node.SearchRef(p.Ref).Value()
	}
}

func (p *Page) Body() (string, error) {
	u, err := p.Url()
	if err != nil {
		return "", err
	}
	if v, e := p.Node.Creeper.CacheGet(u); e {
		return v, nil
	}
	res, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	p.Node.Creeper.CacheSet(u, sb)
	return sb, nil
}

func (p *Page) IsDynamic() bool {
	return p.Town != nil
}

func ParsePage(n *Node, s string) *Page {
	page := new(Page)
	page.Node = n
	page.Raw = s
	if s[0] == '_' {
		page.Ref = s
	} else {
		page.Town = ParseTownLine(s)
		page.Town.Node = n
	}
	page.Index = -1
	return page
}

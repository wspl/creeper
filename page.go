package creeper

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type Page struct {
	Raw string
	Node *Node

	Town *Town
	Ref string

	Index int
}

func (p *Page) Inc() {
	p.Index++
}

func (p *Page) Url() string {
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
				p.Town.Set(k, n.Value())
			}
		}
		return p.Town.Value()
	} else {
		return p.Node.SearchRef(p.Ref).Value()
	}
}

func (p *Page) Body() string {
	if v, e := p.Node.Creeper.Cache_Get(p.Url()); e {
		return v
	}
	u := p.Url()
	res, _ := http.Get(u)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
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
package creeper

import "io/ioutil"

func Open(path string) *Creeper {
	buf, _ := ioutil.ReadFile(path)
	raw := string(buf)
	return New(raw)
}

func New(raw string) *Creeper {
	f := Formatting(raw)
	return NewByFormatted(f)
}

func NewByFormatted(f *Formatted) *Creeper {
	c := new(Creeper)

	c.Nodes = f.Nodes
	for _, n := range c.Nodes {
		n.Creeper = c
	}

	c.Towns = f.Towns
	for _, t := range c.Towns {
		t.Creeper = c
	}

	cache := map[string]string{}
	c.CacheGet = func(k string) (string, bool) {
		v, e := cache[k]
		return v, e
	}
	c.CacheSet = func(k string, v string) {
		cache[k] = v
	}

	return c
}

type Creeper struct {
	Nodes []*Node
	Towns []*Town

	CacheGet func(string) (string, bool)
	CacheSet func(string, string)

	Node *Node
}

func (c *Creeper) Array(key string) *Creeper {
	if c.Node == nil {
		c.Node = c.Nodes[0].SearchFlatScope(key)
	} else {
		c.Node = c.Node.FirstChildNode.SearchFlatScope(key)
	}
	return c
}

func (c *Creeper) String(key string) string {
	v, _ := c.StringE(key)
	return v
}

func (c *Creeper) StringE(key string) (string, error) {
	return c.Node.FirstChildNode.SearchFlatScope(key).Value()
}

func (c *Creeper) Each(cle func(*Creeper)) {
	fstNv := []string{}
	repStor := []string{}
each:
	for {
		v, err := c.Node.Primary().Value()
		v = MD5(v)
		if err != nil {
			continue
		}

		if c.Node.Index == 0 {
			for _, s := range fstNv {
				if s == v {
					break each
				}
			}
			fstNv = append(fstNv, v)
		}

		for _, s := range repStor {
			if s == v {
				continue each
			}
		}
		repStor = append(repStor, v)

		cle(c)
		c.Next()
	}
}

func (c *Creeper) Next() *Creeper {
	c.Node.Inc()
	return c
}

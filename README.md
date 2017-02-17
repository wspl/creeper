![Creeper](https://raw.githubusercontent.com/wspl/creeper/master/art/Creeper.png)

## About

Creeper is a *next-generation* crawler which fetches web page by creeper script. As a cross-platform embedded crawler, you can use it for your news app, subscribe program, etc.

## Get Started

#### Installation

```
$ go get github.com/wspl/creeper
```

#### Hello World!

Create `example.crr`

```
page(@page=1) = "https://news.ycombinator.com/news?p={@page}"

news[]: page -> $("tr.athing")
	title: $(".title a.storylink").text
	site: $(".title span.sitestr").text
```

Then, create `main.go`

```go
package main

import (
	"io/ioutil"
	"github.com/wspl/creeper"
)

func main() {
	buf, _ := ioutil.ReadFile("./example.crr")
	raw := string(buf)
	c := creeper.Creeper_New(raw)
	c.Array("news").Each(func(c *creeper.Creeper) {
		println("===")
		println("title: ", c.String("title"))
		println("site: ", c.String("site"))
	})
}
```

Build and run. Console will print something like:

```
title:  Samsung chief Lee arrested as S.Korean corruption probe deepens
site:  reuters.com
===
title:  ReactOS 0.4.4 Released
site:  reactos.org
===
title:  YC Annual Letter 2017
site:  samaltman.com
===
title:  Let’s not demonize driving, just stop subsidizing it
site:  cityobservatory.org
```
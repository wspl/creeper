![Creeper](https://raw.githubusercontent.com/wspl/creeper/master/art/Creeper.png)
## About

Creeper is a *next-generation* crawler which fetches web page by creeper script. As a cross-platform embedded crawler, you can use it for your news app, subscribe program, etc.

**Warning:** At present this project is still under stage-1 development, please do not use in the production environment.

## Get Started

#### Installation

```
$ go get github.com/wspl/creeper
```

#### Hello World!

Create `hacker_news.crs`

```
page(@page=1) = "https://news.ycombinator.com/news?p={@page}"

news[]: page -> $("tr.athing")
	title: $(".title a.storylink").text
	site: $(".title span.sitestr").text
	link: $(".title a.storylink").href
```

Then, create `main.go`

```go
package main

import "github.com/wspl/creeper"

func main() {
	c := creeper.Open("./hacker_news.crs")
	c.Array("news").Each(func(c *creeper.Creeper) {
		println("title: ", c.String("title"))
		println("site: ", c.String("site"))
		println("link: ", c.String("link"))
		println("===")
	})
}
```

Build and run. Console will print something like:

```
title:  Samsung chief Lee arrested as S.Korean corruption probe deepens
site:  reuters.com
link:  http://www.reuters.com/article/us-southkorea-politics-samsung-group-idUSKBN15V2RD
===
title:  ReactOS 0.4.4 Released
site:  reactos.org
link:  https://reactos.org/project-news/reactos-044-released
===
title:  FeFETs: How this new memory stacks up against existing non-volatile memory
site:  semiengineering.com
link:  http://semiengineering.com/what-are-fefets/
```

## Script Spec

### Town

Town is a lambda like expression for saving (in)mutable string. Most of the time, we used it to store url.

```
page(@page=1, ext) = "https://news.ycombinator.com/news?p={@page}&ext={ext}"
```

When you need town, use it as if you were calling a function:

```
news[]: page(ext="Hello World!") -> $("tr.athing")
```

Hey, you might have noticed that the `@page` parameter is not used. Yeah, it is a special parameter.

Expression in town definition line like `name="something"`, represents parameter `name` has a default value `"something"`.

Incidentally, `@page` is a parameter that will automatically increasing when current page has no more content.


### Node

Nodes is a tree structure that represents the data structure you are going to crawl.

```
news[]: page -> $("tr.athing")
	title: $(".title a.storylink").text
	site: $(".title span.sitestr").text
	link: $(".title a.storylink").href
```

Like `yaml`, nodes distinguishes the hierarchy by indentation.

#### Node Name

Node has name. `title` is a field name, represents a general string data. `news[]` is a array name, represents a parent structure with multiple sub-data.

#### Page

Page indicates where to fetching the field data. It can be a town expression or field reference.

Field reference is a advanced usage of Node, you can found the details in [./eh.crs](./eh.crs).

If a node owned page and fun at the same time, page should on the left of `->`, fun should on the right of `->`. Which is `page -> fun`

#### Fun

Fun represents the data processing process.

There are all supported funs:

| Name      | Parameters                       | Description                              |
| --------- | -------------------------------- | ---------------------------------------- |
| $         | (selector: string)               | CSS selector                             |
| html      |                                  | inner HTML                               |
| text      |                                  | inner text                               |
| outerHTML |                                  | outer HTML                               |
| attr      | (attr: string)                   | attribute value                          |
| style     |                                  | style attribute value                    |
| href      |                                  | href attribute value                     |
| src       |                                  | src attribute value                      |
| calc      | (prec: int)                      | calculate arithmetic expression          |
| match     | (regexp: string)                 | match first sub-string via regular expression |
| expand    | (regexp: string, target: string) | expand matched strings to target string  |



## Author

Plutonist

> [impl.moe](impl.moe) · Github [@wspl](impl.moe) 



## License

```
Copyright (c) 2017 Plutonist
All rights reserved.

Redistribution and use in source and binary forms are permitted
provided that the above copyright notice and this paragraph are
duplicated in all such forms and that any documentation,
advertising materials, and other materials related to such
distribution and use acknowledge that the software was developed
by the Plutonist. The name of the
Plutonist may not be used to endorse or promote products derived
from this software without specific prior written permission.
THIS SOFTWARE IS PROVIDED ``AS IS'' AND WITHOUT ANY EXPRESS OR
IMPLIED WARRANTIES, INCLUDING, WITHOUT LIMITATION, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.
```
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat)](https://opensource.org/licenses/Apache-2.0)
[![PyPI](https://img.shields.io/pypi/status/Django.svg?style=flat)]()
![Creeper](https://raw.githubusercontent.com/wspl/creeper/master/art/Creeper.png)
## About

Creeper is a *next-generation* crawler which fetches web page by creeper script. As a cross-platform embedded crawler, you can use it for your news app, subscribe program, etc.

**Warning:** At present this project is still under stage-1 development, please do not use in the production environment.

## Roadmap

#### Stage-I

* ~~Basic crawling function~~ ← Complete!
* Perfect the code to make it robust
* Add unit test

#### Stage-II

- Decoupling http client and goquery, and add custom interfaces.
- Cross-platform support. Provide Android and iOS SDK.

#### Stage-III

- Add restful api for remote calling
- Add debug api
- Create an IDE based on web for creeper script

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

You might have noticed that the `@page` parameter is not used. Yeah, it is a special parameter.

Expression in town definition line like `name="something"`, represents parameter `name` has a default value `"something"`.

Incidentally, `@page` is a parameter that will automatically increasing when current page has no more content.


### Node

Nodes are tree structure that represent the data structure you are going to crawl.

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

> [impl.moe](https://impl.moe) · Github [@wspl](https://github.com/wspl) 
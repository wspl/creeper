package main

import (
	"github.com/wspl/creeper"
)

func main() {
	c := creeper.Open("./eh.crs")
	c.Array("gallery").Each(func(c *creeper.Creeper) {
		println("title: ", c.MString("title"))
		println("cover: ", c.MString("cover"))
		println("=================================")
	})
}

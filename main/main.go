package main

import (
	"io/ioutil"
	. "github.com/wspl/creeper"
)

func main() {
	buf, _ := ioutil.ReadFile("./example.crr")
	raw := string(buf)
	c := Creeper_New(raw)
	c.Array("gallery").Each(func(c *Creeper) {
		println("title: ", c.String("title"))
		println("cover: ", c.String("cover"))
		println("=================================")
	})
}
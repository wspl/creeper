package main

import (
	"github.com/wspl/creeper"
)

func main() {
	//buf, _ := ioutil.ReadFile("./eh.crr")
	//raw := string(buf)
	//c := New(raw)
	//c.Array("gallery").Each(func(c *Creeper) {
	//	println("title: ", c.String("title"))
	//	println("cover: ", c.String("cover"))
	//	println("=================================")
	//})
	c := creeper.Open("./example.crs")
	c.Array("news").Each(func(c *creeper.Creeper) {
		println("title: ", c.String("title"))
		println("site: ", c.String("site"))
		println("link: ", c.String("link"))
		println("===")
	})
}
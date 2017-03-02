package main

import (
	"github.com/wspl/creeper"
)

func main() {
	c := creeper.Open("./eh.crs")
	c.Array("gallery").Each(func(c *creeper.Creeper) {
		//println("title: ", c.String("title"))
		println("cover: ", c.String("cover"))
		println("===")
	})
	//c = c.Array("gallery")
	//c.Array("pictures").Each(func(c *creeper.Creeper) {
	//	println(c.String("src"))
	//})
	//c := creeper.Open("./example.crs")
	//c.Array("news").Each(func(c *creeper.Creeper) {
	//	println("title: ", c.String("title"))
	//	println("site: ", c.String("site"))
	//	println("link: ", c.String("link"))
	//	println("===")
	//})
}

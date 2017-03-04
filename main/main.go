package main

import (
	"github.com/wspl/creeper"
)

func main() {
	c := creeper.Open("./eh.crs")
	c.Array("gallery").Array("pictures").Each(func(c *creeper.Creeper) {
		//println("title: ", c.String("title"))
		println("t: ", c.String("src"))
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

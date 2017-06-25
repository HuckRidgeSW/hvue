package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

func main() {
	app3 := hvue.NewVM(
		hvue.El("#app-3"),
		hvue.Data("seen", true))
	js.Global.Set("app3", app3)
	// In the JS console, try setting app3.seen to false.
}

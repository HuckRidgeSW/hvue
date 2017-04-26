package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

func main() {
	app3 := hvue.NewVM(
		hvue.El("#app-3"),
		hvue.Data("seen", true))
	js.Global.Set("app3", app3)
}

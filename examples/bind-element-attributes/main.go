package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

func main() {
	app2 := hvue.NewVM(
		hvue.El("#app-2"),
		hvue.Data("message", "You loaded this page on "+time.Now().String()))
	js.Global.Set("app2", app2)
}

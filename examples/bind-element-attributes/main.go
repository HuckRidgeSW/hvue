package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

func main() {
	app2 := hvue.NewVM(
		hvue.El("#app-2"),
		hvue.Data("message", "You loaded this page on "+time.Now().String()))
	js.Global.Set("app2", app2)
	// In the JS console, try setting app2.message to something else.
}

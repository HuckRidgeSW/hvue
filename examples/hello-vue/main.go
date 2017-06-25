package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

func main() {
	app := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"))
	js.Global.Set("app", app)
	// In the JS console, try setting app.message to something else.
}

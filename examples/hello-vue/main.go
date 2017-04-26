package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

func main() {
	app := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"))
	js.Global.Set("app", app)
}

package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

func main() {
	app4 := hvue.NewVM(
		hvue.El("#app-4"),
		hvue.Data("todos", []struct{ Text string }{
			{Text: "Learn JavaScript"},
			{Text: "Learn Vue"},
			{Text: "Build something awesome"}}))
	js.Global.Set("app4", app4)
	// In the JS console, try app4.todos.push({ Text: 'New item' }).
}

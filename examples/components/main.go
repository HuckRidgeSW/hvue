package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

var O = func() *js.Object { return js.Global.Get("Object").New() }

// Several examples in one, from
// https://vuejs.org/v2/guide/components.html

func main() {
	hvue.NewComponent("my-component",
		hvue.Template(`<div>A custom component!</div>`))
	hvue.NewVM(
		hvue.El("#example"))

	//

	type ChildT struct {
		*js.Object
		Template string `js:"template"`
	}
	var Child = hvue.Construct(
		&ChildT{Template: `<div>A custom component, example 2!</div>`},
	)
	hvue.NewVM(
		hvue.El("#example-2"),
		hvue.Component("my-component", Child))
}

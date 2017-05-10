package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

var O = func() *js.Object { return js.Global.Get("Object").New() }

// Several examples in one, from
// https://vuejs.org/v2/guide/components.html

func main() {
	go aRegularComponent()
	go localRegistration()
	go dataMustBeAFunction()
}

func aRegularComponent() {
	hvue.NewComponent("my-component",
		hvue.Template(`<div>A custom component!</div>`))
	hvue.NewVM(
		hvue.El("#example"))
}

func localRegistration() {
	// Local registration
	type ChildT struct {
		*js.Object
		Template string `js:"template"`
	}
	var Child = hvue.NewT(
		&ChildT{Template: `<div>A custom component, example 2!</div>`},
	)
	hvue.NewVM(
		hvue.El("#example-a"),
		hvue.Component("my-component", Child))
}

func dataMustBeAFunction() {
	type DataT struct {
		*js.Object
		counter int `js:"counter"`
	}

	// How NOT to do it: Since all three component instances share the same
	// data object, incrementing one counter increments them all! Ouch.
	data := hvue.NewT(&DataT{counter: 0}).(*DataT)
	hvue.NewComponent(
		"simple-counter1",
		hvue.Template(`<button v-on:click="counter += 1">{{ counter }}</button>`),
		// Return the same object reference for each component instance.  This
		// is an example of how NOT to do data in components.  See the Vue
		// example.
		hvue.DataFunc(func(*hvue.VM) interface{} {
			return data
		}))
	hvue.NewVM(hvue.El("#example-2-a"))

	// Let’s fix this by instead returning a fresh data object:
	hvue.NewComponent(
		"simple-counter2",
		hvue.Template(`<button v-on:click="counter += 1">{{ counter }}</button>`),
		// Return a different object for each component
		hvue.DataFunc(func(*hvue.VM) interface{} {
			return hvue.NewT(&DataT{counter: 0}).(*DataT)
		}))
	hvue.NewVM(hvue.El("#example-2-b"))
}

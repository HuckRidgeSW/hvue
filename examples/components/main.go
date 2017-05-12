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
	go passDataWithProps()
	go propValidation()
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

// https://vuejs.org/v2/guide/components.html#Passing-Data-with-Props
func passDataWithProps() {
	hvue.NewComponent("child",
		hvue.Props("message"),
		hvue.Template(`<span>{{ message }}</span>`))
	hvue.NewVM(hvue.El("#example-3"))
}

// https://vuejs.org/v2/guide/components.html#Prop-Validation
func propValidation() {
	/*
	   Vue.component('example', {
	     props: {
	       // basic type check (`null` means accept any type)
	       propA: Number,
	       // multiple possible types
	       propB: [String, Number],
	       // a required string
	       propC: {
	         type: String,
	         required: true
	       },
	       // a number with default value
	       propD: {
	         type: Number,
	         default: 100
	       },
	       // object/array defaults should be returned from a
	       // factory function
	       propE: {
	         type: Object,
	         default: function () {
	           return { message: 'hello' }
	         }
	       },
	       // custom validator function
	       propF: {
	         validator: function (value) {
	           return value > 10
	         }
	       }
	     }
	   })
	*/
	hvue.NewComponent("child2",
		hvue.Template(`
		<div>propA: {{ propA }}</div>
		`),
		hvue.PropObj("propA",
			hvue.Types(hvue.PNumber)))
	hvue.NewVM(
		hvue.El("#example-4"),
	)
}

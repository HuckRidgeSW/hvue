package main

import (
	"strings"

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
	go counterEvent()
	go counterEventWithChannel()
	go currencyInput()
}

/////////////////////////////////////////////////////////////////////////

func aRegularComponent() {
	hvue.NewComponent("my-component",
		hvue.Template(`<div>A custom component!</div>`))
	hvue.NewVM(
		hvue.El("#example"))
}

/////////////////////////////////////////////////////////////////////////

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

/////////////////////////////////////////////////////////////////////////

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
			// You *can* do the type-assert to its actual type, but you don't
			// *have* to.
			return hvue.NewT(&DataT{counter: 0}).(*DataT)
		}))
	hvue.NewVM(hvue.El("#example-2-b"))
}

/////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Passing-Data-with-Props
func passDataWithProps() {
	hvue.NewComponent("child",
		hvue.Props("message"),
		hvue.Template(`<span>{{ message }}</span>`))
	hvue.NewVM(hvue.El("#example-3"))
}

/////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Prop-Validation
func propValidation() {
	hvue.NewComponent("child2",
		hvue.Template(`
		<div>
			<div>propA: {{ propA }}</div>
			<div>propB: {{ propB }}</div>
			<div>propC: {{ propC }}</div>
			<div>propD: {{ propD }}</div>
			<div>propE: {{ propE.message }}</div>
			<div>propF: {{ propF }}</div>
		</div>
		`),
		// Note kebab-case in the HTML: "prop-a" and so on.
		hvue.PropObj("propA",
			hvue.Types(hvue.PNumber)),
		hvue.PropObj("propB",
			hvue.Types(hvue.PString, hvue.PNumber)),
		hvue.PropObj("propC",
			hvue.Types(hvue.PString), hvue.Required),
		hvue.PropObj("propD",
			hvue.Types(hvue.PNumber),
			hvue.Default(100)),
		hvue.PropObj("propE",
			hvue.Types(hvue.PObject),
			hvue.DefaultFunc(func(*hvue.VM) interface{} {
				return js.M{"message": "hello"}
			})),
		hvue.PropObj("propF",
			hvue.Validator(func(vm *hvue.VM, value *js.Object) interface{} {
				return value.Int() > 10
			})),
	)
	hvue.NewVM(hvue.El("#example-4"))
}

/////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Using-v-on-with-Custom-Events

type ButtonCounterT struct {
	*js.Object
	Counter int `js:"counter"`
}

type CounterEventT struct {
	*js.Object
	Total int `js:"total"`
}

func counterEvent() {
	hvue.NewComponent("button-counter",
		hvue.Template(`<button v-on:click="Increment">{{ counter }}</button>`),
		hvue.DataFunc(func(*hvue.VM) interface{} {
			return hvue.NewT(&ButtonCounterT{Counter: 0})
		}),
		hvue.MethodsOf(&ButtonCounterT{}))
	hvue.NewVM(
		hvue.El("#counter-event-example"),
		hvue.DataS(hvue.NewT(&CounterEventT{Total: 0})),
		hvue.MethodsOf(&CounterEventT{}))
}

func (o *ButtonCounterT) Increment(vm *hvue.VM) {
	o.Counter++
	vm.Emit("increment")

}

func (o *CounterEventT) IncrementTotal(vm *hvue.VM) {
	o.Total++
}

/////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Using-v-on-with-Custom-Events
// again, but with a channel.

// Note that this is a proof of concept, to show one way to use a channel.
// This particular way isn't very good, as it just re-implements part of the
// Vue event model, badly.

type ButtonCounterWithChannelT struct {
	*js.Object
	Counter int `js:"counter"`
	eventCh chan string
}

// Reused from above
// type CounterEventT struct {
// 	*js.Object
// 	Total int `js:"total"`
// }

func counterEventWithChannel() {
	eventCh := make(chan string)
	hvue.NewComponent("button-counter-with-channel",
		hvue.Template(`<button v-on:click="Increment">{{ counter }}</button>`),
		hvue.DataFunc(func(*hvue.VM) interface{} {
			return hvue.NewT(&ButtonCounterWithChannelT{
				Counter: 0,
				eventCh: eventCh,
			})
		}),
		hvue.MethodsOf(&ButtonCounterWithChannelT{}))

	data := hvue.NewT(&CounterEventT{Total: 0}).(*CounterEventT)
	vm := hvue.NewVM(
		hvue.El("#counter-event-example-with-channel"),
		hvue.DataS(data),
		hvue.MethodsOf(&CounterEventT{}))

	go func() {
		for event := range eventCh {
			switch event {
			case "increment":
				data.IncrementTotal(vm)
			}
		}
	}()
}

func (o *ButtonCounterWithChannelT) Increment(vm *hvue.VM) {
	o.Counter++
	o.eventCh <- "increment"

}

// Reused from above
// func (o *CounterEventT) IncrementTotal(vm *hvue.VM) {
// 	o.Total++
// }

/////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Form-Input-Components-using-Custom-Events

type CurrencyData struct {
	*js.Object
	Price string `js:"price"`
}

func currencyInput() {
	hvue.NewComponent("currency-input",
		hvue.Template(`
		<span>
		  $
		  <input
		    ref="input"
		    v-bind:value="value"
		    v-on:input="UpdateValue($event.target.value)">
		</span>
		`),
		hvue.Props("value"),
		hvue.DataFunc(func(*hvue.VM) interface{} {
			return hvue.NewT(&CurrencyData{})
		}),
		// hvue.MethodsOf(&CurrencyData{}),
		hvue.Method("UpdateValue", func(vm *hvue.VM, value string) {
			// Remove whitespace on either side
			formattedValue := strings.TrimSpace(value)
			formattedValue = formattedValue[:dotPlus3(formattedValue)]
			// If the value was not already normalized,
			// manually override it to conform
			if formattedValue != value {
				vm.Refs("input").Set("value", formattedValue)
			}
			vm.Emit("input", js.Global.Get("Number").Invoke(formattedValue))
		}),
	)
	hvue.NewVM(
		hvue.El("#currency-input-example"),
		hvue.DataS(hvue.NewT(&CurrencyData{Price: ""})))
}

// Instead of updating the value directly, this
// method is used to format and place constraints
// on the input's value
func (_ *CurrencyData) UpdateValue(vm *hvue.VM, value string) {
	// Remove whitespace on either side
	formattedValue := strings.TrimSpace(value)
	formattedValue = formattedValue[:dotPlus3(formattedValue)]
	// If the value was not already normalized,
	// manually override it to conform
	if formattedValue != value {
		vm.Refs("input").Set("value", formattedValue)
	}
	vm.Emit("input", js.Global.Get("Number").Invoke(formattedValue))
}

func dotPlus3(value string) int {
	if i := strings.Index(value, "."); i == -1 || i+3 > len(value) {
		return len(value)
	} else {
		return i + 3
	}
}

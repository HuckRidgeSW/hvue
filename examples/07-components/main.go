package main

// From https://vuejs.org/v2/guide/components.html

import (
	"strconv"
	"strings"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"

	"github.com/huckridgesw/hvue"
)

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
	go moreRobustCurrencyInput()

	select {}
}

////////////////////////////////////////////////////////////////////////////////

func aRegularComponent() {
	hvue.NewComponent("my-component",
		hvue.Template(`<div>A custom component!</div>`))
	hvue.NewVM(
		hvue.El("#example"))
}

////////////////////////////////////////////////////////////////////////////////

func localRegistration() {
	hvue.NewVM(
		hvue.El("#example-a"),
		hvue.Component("my-local-component",
			hvue.Template(`<div>A custom component, example 2!</div>`)))
}

////////////////////////////////////////////////////////////////////////////////

func dataMustBeAFunction() {
	type DataT struct {
		js.Value
	}

	// How NOT to do it: Since all three component instances share the same
	// data object, incrementing one counter increments them all!  Ouch.
	data := hvue.Map2Obj(hvue.M{"counter": 0})
	hvue.NewComponent(
		"simple-counter1",
		hvue.Template(`<button v-on:click="counter += 1">{{ counter }}</button>`),
		// Return the same object reference for each component instance.  This
		// is an example of how NOT to do data in components.  See the Vue
		// example.
		//
		// Have to use a custom ComponentOption function, because hvue.DataFunc
		// actually makes it impossible to not return a new object each time.
		func(c *hvue.Config) {
			c.DataType = js.TypeFunction
			c.Set("data", js.NewCallback(
				func(js.Value, []js.Value) interface{} {
					// Return the same object each time (don't do this).
					return data
				}))
		})
	hvue.NewVM(hvue.El("#example-2-a"))

	// Let’s fix this by instead returning a fresh data object:
	hvue.NewComponent(
		"simple-counter2",
		hvue.Template(`<button v-on:click="counter += 1">{{ counter }}</button>`),
		// Return a different object for each component
		hvue.DataFunc(func(_ *hvue.VM, o js.Value) interface{} {
			o.Set("counter", 0)
			return &DataT{Value: o}
		}),
	)
	hvue.NewVM(hvue.El("#example-2-b"))
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Passing-Data-with-Props
func passDataWithProps() {
	hvue.NewComponent("child",
		hvue.Props("message"),
		hvue.Template(`<span>{{ message }}</span>`))
	hvue.NewVM(hvue.El("#example-3"))
}

////////////////////////////////////////////////////////////////////////////////

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
			hvue.DefaultFunc(hvue.Map2Obj(hvue.M{"message": "hello"}))),
		hvue.PropObj("propF",
			hvue.Validator(func(_ *hvue.VM, value js.Value) interface{} {
				return value.Int() > 10
			})),
	)
	hvue.NewVM(hvue.El("#example-4"))
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Using-v-on-with-Custom-Events

type ButtonCounterT struct {
	js.Value
}

func (b *ButtonCounterT) SetCounter(new int) { b.Set("counter", new) }
func (o *ButtonCounterT) Increment(vm *hvue.VM) {
	o.Set("counter", o.Get("counter").Int()+1)
	vm.Emit("increment")
}

type CounterEventT struct {
	js.Value
}

func (b *CounterEventT) Total() int              { return b.Get("total").Int() }
func (b *CounterEventT) SetTotal(new int)        { b.Set("total", new) }
func (o *CounterEventT) IncrementTotal(*hvue.VM) { o.SetTotal(o.Total() + 1) }

func counterEvent() {
	hvue.NewComponent("button-counter",
		hvue.Template(`<button v-on:click="Increment">{{ counter }}</button>`),
		hvue.DataFunc(func(_ *hvue.VM, o js.Value) interface{} {
			data := &ButtonCounterT{Value: o}
			data.SetCounter(0)
			return data
		}),
		hvue.MethodsOf(&ButtonCounterT{}))
	data2 := &CounterEventT{Value: hvue.Map2Obj(hvue.M{"total": 0})}
	hvue.NewVM(
		hvue.El("#counter-event-example"),
		hvue.DataS(data2, data2.Value),
		hvue.MethodsOf(&CounterEventT{}))
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Using-v-on-with-Custom-Events
// again, but with a channel.

// Note that this is a proof of concept, to show one way to use a channel.
// This particular way isn't very good, as it just re-implements part of the
// Vue event model, badly.

type ButtonCounterWithChannelT struct {
	js.Value
	eventCh chan string
}

func (b *ButtonCounterWithChannelT) Counter() int       { return b.Get("counter").Int() }
func (b *ButtonCounterWithChannelT) SetCounter(new int) { b.Set("counter", new) }

func counterEventWithChannel() {
	eventCh := make(chan string)
	hvue.NewComponent("button-counter-with-channel",
		hvue.Template(`<button v-on:click="Increment">{{ counter }}</button>`),
		hvue.DataFunc(func(_ *hvue.VM, o js.Value) interface{} {
			data := &ButtonCounterWithChannelT{
				Value:   o,
				eventCh: eventCh,
			}
			data.SetCounter(0)
			return data
		}),
		hvue.MethodsOf(&ButtonCounterWithChannelT{}))

	data := &CounterEventT{Value: hvue.Map2Obj(hvue.M{"total": 0})}
	vm := hvue.NewVM(
		hvue.El("#counter-event-example-with-channel"),
		hvue.DataS(data, data.Value),
		hvue.MethodsOf(data))

	go func() {
		for event := range eventCh {
			switch event {
			case "increment":
				data.IncrementTotal(vm)
			}
		}
	}()
}

func (o *ButtonCounterWithChannelT) Increment(*hvue.VM) {
	o.SetCounter(o.Counter() + 1)
	o.eventCh <- "increment"
}

// CounterEventT and its method IncrementTotal reused from above.

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/components.html#Form-Input-Components-using-Custom-Events

type CurrencyData struct {
	js.Value
}

func (c *CurrencyData) Price() string       { return c.Get("price").String() }
func (c *CurrencyData) SetPrice(new string) { c.Set("price", new) }

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
		hvue.DataFunc(func(_ *hvue.VM, o js.Value) interface{} {
			data := &CurrencyData{Value: o}
			data.SetPrice("0")
			return data
		}),

		// Show two ways of adding the UpdateValue method:

		// #1: Automatically add all methods of *CurrencyData.  NOTE: that
		// UpdateValue is not actually a method of CurrencyData, this is just an
		// example call of hvue.MethodsOf.
		// hvue.MethodsOf(&CurrencyData{}),

		// #2: Add this closure as a single named method:
		hvue.Method("UpdateValue", func(vm *hvue.VM, value string) {
			// Remove whitespace on either side
			formattedValue := strings.TrimSpace(value)
			formattedValue = formattedValue[:dotPlus3(formattedValue)]
			// If the value was not already normalized,
			// manually override it to conform
			if formattedValue != value {
				vm.Refs("input").Set("value", formattedValue)
			}
			vm.Emit("input", js.Global().Get("Number").Invoke(formattedValue))
		}),
	)
	data := &CurrencyData{Value: hvue.Map2Obj(hvue.M{"price": ""})}
	hvue.NewVM(
		hvue.El("#currency-input-example"),
		hvue.DataS(data, data.Value))
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
	vm.Emit("input", js.Global().Get("Number").Invoke(formattedValue))
}

func dotPlus3(value string) int {
	if i := strings.Index(value, "."); i == -1 || i+3 > len(value) {
		return len(value)
	} else {
		return i + 3
	}
}

////////////////////////////////////////////////////////////////////////////////

// Here’s a more robust currency filter
// Still from https://vuejs.org/v2/guide/components.html#Form-Input-Components-using-Custom-Events

type CurrencyInputT struct {
	js.Value
}

func (c *CurrencyInputT) Price() float64      { return c.Get("price").Float() }
func (c *CurrencyInputT) Shipping() float64   { return c.Get("shipping").Float() }
func (c *CurrencyInputT) Handling() float64   { return c.Get("handling").Float() }
func (c *CurrencyInputT) Discount() float64   { return c.Get("discount").Float() }
func (c *CurrencyInputT) SetTotal(new string) { c.Set("total", new) }

func moreRobustCurrencyInput() {
	hvue.NewComponent("currency-input2",
		hvue.Template(`
        <div>
          <label v-if="label">{{ label }}</label>
          $
          <input
            ref="input"
            v-bind:value="value"
            v-on:input="UpdateValue($event.target.value)"
				v-on:focus="SelectAll"
            v-on:blur="FormatValue"
          >
        </div>`),
		hvue.PropObj(
			"value",
			hvue.Types(hvue.PNumber),
			hvue.Default(0)),
		hvue.PropObj(
			"label",
			hvue.Types(hvue.PString),
			hvue.Default("")),
		hvue.Mounted(func(vm *hvue.VM) {
			vm.Call("FormatValue")
		}),
		hvue.MethodsOf(&CurrencyInputT{}),
	)

	data := &CurrencyInputT{
		Value: hvue.Map2Obj(hvue.M{
			"price":    0,
			"shipping": 0,
			"handling": 0,
			"discount": 0,
			"total":    "",
		})}
	watch := func(*hvue.VM) {
		data.SetTotal(strconv.FormatFloat((data.Price()*100+
			data.Shipping()*100+
			data.Handling()*100-
			data.Discount()*100)/100, 'f', 2, 32))
	}
	hvue.NewVM(
		hvue.El("#app"),
		hvue.DataS(data, data.Value),
		hvue.Watch("price", watch),
		hvue.Watch("shipping", watch),
		hvue.Watch("handling", watch),
		hvue.Watch("discount", watch),
	)
}

func (c *CurrencyInputT) UpdateValue(vm *hvue.VM, value js.Value) {
	result := js.Global().Get("currencyValidator").
		Call("parse", value, vm.Get("value"))
	if result.Get("warning") != js.Undefined() {
		vm.Refs("input").Set("value", result.Get("value"))
	}
	vm.Emit("input", result.Get("value"))
}

func (c *CurrencyInputT) FormatValue(vm *hvue.VM) {
	vm.Refs("input").Set("value",
		js.Global().Get("currencyValidator").Call("format", vm.Get("value")))
}

func (c *CurrencyInputT) SelectAll(_ *hvue.VM, event *hvue.Event) {
	event.Target().Select()
}

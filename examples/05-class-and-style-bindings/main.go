package main

// Several examples in one, from
// https://vuejs.org/v2/guide/class-and-style.html.

import (
	"time"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"

	"github.com/huckridgesw/hvue"
)

type ClassObject struct {
	js.Value
}

func main() {
	go dynamically_toggle_classes()
	go doesnt_have_to_be_inline()
	go bind_to_a_computed_property()
	go binding_styles()

	select {}
}

////////////////////////////////////////////////////////////////////////////////

func dynamically_toggle_classes() {
	// "We can pass an object to v-bind:class to dynamically toggle classes"
	data1 := &struct {
		js.Value
	}{Value: hvue.NewObject()}
	data1.Set("isActive", false)

	app1 := hvue.NewVM(
		hvue.El("#object-syntax-1"),
		hvue.DataS(data1, data1.Value))
	js.Global().Set("app1", app1.Value)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			data1.Set("isActive", !data1.Get("isActive").Bool())
		}
	}()
}

////////////////////////////////////////////////////////////////////////////////

func doesnt_have_to_be_inline() {
	// "The bound object doesn’t have to be inline"
	data2 := &struct {
		js.Value
		*ClassObject
	}{Value: hvue.NewObject(), ClassObject: &ClassObject{Value: hvue.NewObject()}}
	data2.Set("classObject", data2.ClassObject.Value)
	data2.ClassObject.Set("active", true)
	data2.ClassObject.Set("text-danger", false)

	app2 := hvue.NewVM(
		hvue.El("#object-syntax-2"),
		hvue.DataS(data2, data2.Value))
	js.Global().Set("app2", app2.Value)
}

////////////////////////////////////////////////////////////////////////////////

func bind_to_a_computed_property() {
	// "We can also bind to a computed property that returns an object"
	data3 := &struct {
		js.Value
	}{Value: hvue.Map2Obj(hvue.M{
		"isActive": false,
		"error":    hvue.NewObject(),
	})}
	js.Global().Set("data3", data3.Value)

	hvue.NewVM(
		hvue.El("#object-syntax-3"),
		hvue.DataS(data3, data3.Value),
		hvue.Computed(
			"classObject",
			func(*hvue.VM) interface{} {
				return hvue.Map2Obj(hvue.M{
					"active": hvue.Truthy(data3.Get("isActive")) &&
						hvue.Falsy(data3.Get("error")),
					"text-danger": hvue.Truthy(data3.Get("error")) &&
						data3.Get("error").Get("type").String() == "fatal",
				})
			}))

	// In the JS console, try
	//
	//   Vue.set(data3.error, "type", "fatal")
	//
	// and watch this example's TextDanger line (the 2nd one down)
}

////////////////////////////////////////////////////////////////////////////////

type StyleObject struct {
	js.Value
}

func (so *StyleObject) SetColor(new string)    { so.Set("color", new) }
func (so *StyleObject) SetFontSize(new string) { so.Set("fontSize", new) }

func binding_styles() {
	// Binding styles
	data4 := &struct {
		js.Value
		*StyleObject
	}{
		Value:       hvue.NewObject(),
		StyleObject: &StyleObject{Value: hvue.NewObject()},
	}
	data4.Set("styleObject", data4.StyleObject.Value)
	data4.SetColor("red")
	data4.SetFontSize("13px")

	hvue.NewVM(
		hvue.El("#object-syntax-4"),
		hvue.DataS(data4, data4.Value))
}

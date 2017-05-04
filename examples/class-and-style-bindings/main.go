package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

// Several examples in one, from
// https://vuejs.org/v2/guide/class-and-style.html.

func main() {
	// "We can pass an object to v-bind:class to dynamically toggle classes"
	data1 := &struct {
		*js.Object
		IsActive bool `js:"isActive"`
	}{Object: hvue.NewObject()}
	data1.IsActive = false

	hvue.NewVM(
		hvue.El("#object-syntax-1"),
		hvue.DataS(data1))

	go func() {
		time.Sleep(time.Second)
		data1.IsActive = true
		println("isActive:", data1.IsActive)
	}()

	// "The bound object doesn’t have to be inline"
	type ClassObject struct {
		*js.Object
		Active     bool `js:"active"`
		TextDanger bool `js:"text_danger"`
	}
	data2 := &struct {
		*js.Object
		*ClassObject `js:"classObject"`
	}{Object: hvue.NewObject()}
	data2.ClassObject = &ClassObject{Object: hvue.NewObject()}
	data2.ClassObject.Active = false
	data2.ClassObject.TextDanger = false

	hvue.NewVM(
		hvue.El("#object-syntax-2"),
		hvue.DataS(data2))

	// "We can also bind to a computed property that returns an object"
	type errorType struct {
		*js.Object
		Type string `js:"type"`
	}
	data3 := &struct {
		*js.Object
		IsActive bool       `js:"isActive"`
		Error    *errorType `js:"error"`
	}{Object: hvue.NewObject()}
	data3.IsActive = false
	data3.Error = nil

	hvue.NewVM(
		hvue.El("#object-syntax-3"),
		hvue.DataS(data3),
		hvue.Computed(
			"classObject",
			func(vm *hvue.VM) interface{} {
				co := &ClassObject{Object: hvue.NewObject()}
				co.Active = data3.IsActive && data3.Error.Object == nil
				co.TextDanger = data3.Error.Object != nil &&
					data3.Error.Type == "fatal"
				return co
			}))

	// Binding styles
	type StyleObject struct {
		*js.Object
		Color    string `js:"color"`
		FontSize string `js:"fontSize"`
	}
	data4 := &struct {
		*js.Object
		*StyleObject `js:"styleObject"`
	}{Object: hvue.NewObject()}
	data4.StyleObject = &StyleObject{Object: hvue.NewObject()}
	// As of this writing, you can't assign data4.Color or FontSize directly.
	data4.StyleObject.Color = "red"
	data4.StyleObject.FontSize = "13px"

	hvue.NewVM(
		hvue.El("#object-syntax-4"),
		hvue.DataS(data4))
}

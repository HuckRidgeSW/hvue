package main

// This file has code from the Vue Guide Introduction page.
// https://github.com/HuckRidgeSW/hvue/tree/master/examples/01-introduction.

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/huckridgesw/hvue"
)

func main() {
	go helloVue()
	go bindElementAttributes()
	go togglePresenceOfElement()
	go displayAList()
	go pressAButton()
	go twoWayBinding()
	go composingWithComponents()

	select {}
}

////////////////////////////////////////////////////////////////////////////////

func helloVue() {
	app := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"))
	js.Global().Set("app", app.Value)
	// In the JS console, try setting app.message to something else.
}

////////////////////////////////////////////////////////////////////////////////

func bindElementAttributes() {
	app2 := hvue.NewVM(
		hvue.El("#app-2"),
		hvue.Data("message", "You loaded this page on "+time.Now().String()))
	js.Global().Set("app2", app2.Value)
	// In the JS console, try setting app2.message to something else.
}

////////////////////////////////////////////////////////////////////////////////

func togglePresenceOfElement() {
	app3 := hvue.NewVM(
		hvue.El("#app-3"),
		hvue.Data("seen", true))
	js.Global().Set("app3", app3.Value)
	// In the JS console, try setting app3.seen to false.
}

////////////////////////////////////////////////////////////////////////////////

func displayAList() {
	data := hvue.NewArray()
	for i, v := range []string{"Learn JavaScript", "Learn Vue", "Build something awesome"} {
		data.SetIndex(i, map[string]interface{}{"Text": v})
	}

	app4 := hvue.NewVM(
		hvue.El("#app-4"),
		hvue.Data("todos", data))
	js.Global().Set("app4", app4.Value)
	// In the JS console, try app4.todos.push({ Text: 'New item' }).
}

////////////////////////////////////////////////////////////////////////////////

type Data5 struct{ js.Value }

func (d *Data5) Message() string       { return d.Get("message").String() }
func (d *Data5) SetMessage(new string) { d.Set("message", new) }

func (d *Data5) ReverseMessage() {
	d.SetMessage(reverse(d.Message()))
}

func pressAButton() {
	d5 := &Data5{Value: hvue.NewObject()}
	d5.SetMessage("Hello, Vue!")
	hvue.NewVM(
		hvue.El("#app-5"),
		hvue.DataS(d5, d5.Value),
		hvue.MethodsOf(&Data5{})) // FIXME: Could &Data5{} be d5?
}

func reverse(s string) string {
	runes := strings.Split(s, "")
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return strings.Join(runes, "")
}

////////////////////////////////////////////////////////////////////////////////

func twoWayBinding() {
	hvue.NewVM(
		hvue.El("#app-6"),
		hvue.Data("message", "Hello Vue!"))
}

////////////////////////////////////////////////////////////////////////////////

type Data7 struct{ js.Value }
type ListItem7 struct{ js.Value }

func composingWithComponents() {
	hvue.NewComponent("todo-item",
		hvue.Props("todo"),
		hvue.Template("<li>{{ todo.text }}</li>"))

	d7 := NewData7(
		"Vegetables",
		"Cheese",
		"Whatever else humans are supposed to eat")

	app7 := hvue.NewVM(
		hvue.El("#app-7"),
		hvue.DataS(d7, d7.Value))
	js.Global().Set("app7", app7.Value)
	// In the JS console, try app7.groceryList.push({text: "a new item"})
}

func NewData7(texts ...string) *Data7 {
	d := &Data7{Value: hvue.NewObject()}
	d.Set("groceryList", hvue.NewArray())
	gl := d.Get("groceryList")
	for i, v := range texts {
		o := hvue.NewObject()
		o.Set("text", v)
		gl.SetIndex(i, o)
	}
	return d
}

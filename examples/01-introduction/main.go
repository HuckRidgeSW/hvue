package main

// This file has code from the Vue Guide Introduction page.
// https://github.com/HuckRidgeSW/hvue/tree/master/examples/01-introduction.

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
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
}

////////////////////////////////////////////////////////////////////////////////

func helloVue() {
	app := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"))
	js.Global.Set("app", app)
	// In the JS console, try setting app.message to something else.
}

////////////////////////////////////////////////////////////////////////////////

func bindElementAttributes() {
	app2 := hvue.NewVM(
		hvue.El("#app-2"),
		hvue.Data("message", "You loaded this page on "+time.Now().String()))
	js.Global.Set("app2", app2)
	// In the JS console, try setting app2.message to something else.
}

////////////////////////////////////////////////////////////////////////////////

func togglePresenceOfElement() {
	app3 := hvue.NewVM(
		hvue.El("#app-3"),
		hvue.Data("seen", true))
	js.Global.Set("app3", app3)
	// In the JS console, try setting app3.seen to false.
}

////////////////////////////////////////////////////////////////////////////////

func displayAList() {
	app4 := hvue.NewVM(
		hvue.El("#app-4"),
		hvue.Data("todos", []struct{ Text string }{
			{Text: "Learn JavaScript"},
			{Text: "Learn Vue"},
			{Text: "Build something awesome"}}))
	js.Global.Set("app4", app4)
	// In the JS console, try app4.todos.push({ Text: 'New item' }).
}

////////////////////////////////////////////////////////////////////////////////

type Data5 struct {
	*js.Object
	Message string `js:"message"`
}

func pressAButton() {
	hvue.NewVM(
		hvue.El("#app-5"),
		hvue.DataS(hvue.NewT(&Data5{Message: "Hello, Vue!"})),
		hvue.MethodsOf(&Data5{}))
}

func (d *Data5) ReverseMessage() {
	d.Message = reverse(d.Message)
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

type Data7 struct {
	*js.Object
	GroceryList []*ListItem7 `js:"groceryList"`
}

type ListItem7 struct {
	*js.Object
	Text string `js:"text"`
}

func composingWithComponents() {
	hvue.NewComponent("todo-item",
		hvue.Props("todo"),
		hvue.Template("<li>{{ todo.text }}</li>"))

	// This compiles and runs but wouldn't actually work in practice.  It'd be
	// nice to write a function that could take this and copy it into a new
	// structure with all the *js.Object slots initialized correctly.
	testData := &Data7{
		GroceryList: []*ListItem7{
			&ListItem7{Text: "stuff"},
		},
	}
	println("testData is", testData)

	hvue.NewVM(
		hvue.El("#app-7"),
		hvue.DataS(NewData(
			"Vegetables",
			"Cheese",
			"Whatever else humans are supposed to eat")))
}

func NewData(texts ...string) *Data7 {
	d := &Data7{Object: hvue.NewObject()}
	d.GroceryList = NewGroceryList(texts)
	return d
}

func NewGroceryList(texts []string) []*ListItem7 {
	list := make([]*ListItem7, len(texts))
	for i, text := range texts {
		list[i] = NewListItem(text)
	}
	return list
}

func NewListItem(Text string) *ListItem7 {
	item := &ListItem7{Object: js.Global.Get("Object").New()}
	item.Text = Text
	return item
}

package main

// From https://vuejs.org/v2/guide/list.html

import (
	"strconv"
	"time"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/huckridgesw/hvue"
)

// Several examples in one, from
// https://vuejs.org/v2/guide/list.html

// ItemT is used in all three examples.
type ItemT struct {
	js.Value
}

func (i *ItemT) Message() string       { return i.Get("message").String() }
func (i *ItemT) SetMessage(new string) { i.Set("message", new) }

func main() {
	go BasicUsage()
	go v_for_block()
	go simpleTodoList()

	select {}
}

////////////////////////////////////////////////////////////////////////////////
// Basic usage

type Data1 struct {
	js.Value
}

func (d *Data1) Items() js.Value              { return d.Get("items") }
func (d *Data1) Item(i int) js.Value          { return d.Items().Index(i) }
func (d *Data1) SetItems(new js.Value)        { d.Set("items", new) }
func (d *Data1) SetItemN(i int, new js.Value) { d.Items().SetIndex(i, new) }

func BasicUsage() {
	// Basic usage
	data := &Data1{Value: hvue.NewObject()}
	data.SetItems(hvue.NewArray())
	data.SetItemN(0, NewItem("Foo"))
	data.SetItemN(1, NewItem("Bar"))

	hvue.NewVM(
		hvue.El("#example-1"),
		hvue.DataS(data, data.Value))

	// Demonstrate reactivity
	time.Sleep(500 * time.Millisecond)
	data.Items().Index(0).Set("message", "Baz")

	time.Sleep(500 * time.Millisecond)
	data.Items().Index(1).Set("message", "Qux")

	// Add a slice element.
	time.Sleep(500 * time.Millisecond)
	data.Items().Call("push", NewItem("Quux"))

	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		data.Items().Call("push", NewItem(strconv.Itoa(i)))
	}
	time.Sleep(500 * time.Millisecond)
	randomSlice := data.Items().Call("slice", 3, 5)
	randomSlice.Index(1).Set("message", "I am randomSlice[1]")

	time.Sleep(500 * time.Millisecond)
	hvue.Push(data.Items(), NewItem("Quuz"))
}

////////////////////////////////////////////////////////////////////////////////
// v-for example

type Data2 struct {
	js.Value
}

func (d *Data2) SetItems(new js.Value)        { d.Set("items", new) }
func (d *Data2) SetParentMessage(new string)  { d.Set("parentMessage", new) }
func (d *Data2) SetItemN(i int, new js.Value) { d.Get("items").SetIndex(i, new) }

func v_for_block() {
	data := &Data2{Value: hvue.NewObject()}
	data.SetParentMessage("Parent")
	data.SetItems(hvue.NewArray())
	data.SetItemN(0, NewItem("Foo"))
	data.SetItemN(1, NewItem("Bar"))

	hvue.NewVM(
		hvue.El("#example-2"),
		hvue.DataS(data, data.Value))
}

// NewItem is used in both of the above examples
func NewItem(m string) js.Value {
	i := &ItemT{Value: hvue.NewObject()}
	i.SetMessage(m)
	return i.Value
}

////////////////////////////////////////////////////////////////////////////////
// Simple todo list example

type Data3 struct {
	js.Value // `js:"newTodoText:string; todos: []string;"`
}

func (d *Data3) NewTodoText() string       { return d.Get("newTodoText").String() }
func (d *Data3) SetNewTodoText(new string) { d.Set("newTodoText", new) }
func (d *Data3) Todos() js.Value           { return d.Get("todos") }
func (d *Data3) SetTodos(new js.Value)     { d.Set("todos", new) }
func (d *Data3) SetTodosFromStrings(new ...string) {
	todo := d.Todos()
	for i, s := range new {
		todo.SetIndex(i, s)
	}
}

// Event handler, called from html.
func (d *Data3) AddNewTodo() {
	d.Todos().Call("push", d.NewTodoText())
	d.SetNewTodoText("")
}

func simpleTodoList() {
	hvue.NewComponent("todo-item",
		hvue.Template(`
			<li>
			  {{ title }}
			  <button v-on:click="$emit('remove')">X</button>
			</li>
		`),
		hvue.Props("title"))

	data := &Data3{Value: hvue.NewObject()}
	data.SetNewTodoText("")
	data.SetTodos(hvue.NewArray())
	data.SetTodosFromStrings(
		"Do the dishes",
		"Take out the trash",
		"Mow the lawn",
	)

	hvue.NewVM(
		hvue.El("#todo-list-example"),
		hvue.DataS(data, data.Value),
		hvue.MethodsOf(&Data3{}))

	// Show how to update an array element in place.
	time.Sleep(500 * time.Millisecond)
	hvue.Set(data.Todos(), 1, "UPDATE: Take out the papers and the trash")
}

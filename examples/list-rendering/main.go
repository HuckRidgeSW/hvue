package main

import (
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

var O = func() *js.Object { return js.Global.Get("Object").New() }

// Several examples in one, from
// https://vuejs.org/v2/guide/list.html

// ItemT is used in all three examples.
type ItemT struct {
	*js.Object
	Message string `js:"message"`
}

func main() {
	go BasicUsage()
	go v_for_block()
	go simpleTodoList()
}

////////////////////////////////////////////////////////
// Basic usage

type Data1 struct {
	*js.Object
	Items []*ItemT `js:"items"`
}

func BasicUsage() {
	// Basic usage
	data := &Data1{Object: O()}
	data.Items = []*ItemT{
		NewItem("Foo"),
		NewItem("Bar"),
	}
	hvue.NewVM(
		hvue.El("#example-1"),
		hvue.DataS(data))

	// Demonstrate reactivity
	time.Sleep(500 * time.Millisecond)
	data.Items[0].Message = "Baz"

	time.Sleep(500 * time.Millisecond)
	data.Items[1].Message = "Qux"

	// Add a slice element.
	time.Sleep(500 * time.Millisecond)
	// I'm surprised that this actually works, but it does appear to.  I tried
	// appending 1000 items, and they're all reactive in the usual way.
	// Whatever magic GopherJS and Vue do appears to work together.
	data.Items = append(data.Items, NewItem("Quux"))

	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		data.Items = append(data.Items, NewItem(strconv.Itoa(i)))
	}
	time.Sleep(500 * time.Millisecond)
	randomSlice := data.Items[3:5]
	randomSlice[1].Message = "I am randomSlice[1]"

	time.Sleep(500 * time.Millisecond)
	// Here's the way I would have expected that you'd have to do the above
	// appending -- slightly longer and more vulnerable to typos (since the
	// field name is just a string):
	hvue.Push(data.Get("items"), NewItem("Quuz"))
}

////////////////////////////////////////////////////////
// v-for example

func v_for_block() {
	type Data2 struct {
		*js.Object
		ParentMessage string   `js:"parentMessage"`
		Items         []*ItemT `js:"items"`
	}

	data := &Data2{Object: O()}
	data.ParentMessage = "Parent"
	data.Items = []*ItemT{
		NewItem("Foo"),
		NewItem("Bar"),
	}

	hvue.NewVM(
		hvue.El("#example-2"),
		hvue.DataS(data))
}

// NewItem is used in both of the above examples
func NewItem(m string) *ItemT {
	i := &ItemT{Object: O()}
	i.Message = m
	return i
}

////////////////////////////////////////////////////////
// Simple todo list example

type Data3 struct {
	*js.Object
	NewTodoText string   `js:"newTodoText"`
	Todos       []string `js:"todos"`
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

	data := &Data3{Object: O()}
	data.NewTodoText = ""
	data.Todos = []string{
		"Do the dishes",
		"Take out the trash",
		"Mow the lawn",
	}

	hvue.NewVM(
		hvue.El("#todo-list-example"),
		hvue.DataS(data),
		hvue.MethodsOf(&Data3{}))
}

func (d *Data3) AddNewTodo() {
	d.Todos = append(d.Todos, d.NewTodoText)
	d.NewTodoText = ""
}

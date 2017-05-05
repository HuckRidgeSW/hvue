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

type ItemT struct {
	*js.Object
	Message string `js:"message"`
}

type DataT struct {
	*js.Object
	Items []*ItemT `js:"items"`
}

func main() {
	// Basic usage
	data1 := &DataT{Object: O()}
	data1.Items = []*ItemT{
		NewItem("Foo"),
		NewItem("Bar"),
	}
	hvue.NewVM(
		hvue.El("#example-1"),
		hvue.DataS(data1))

	// Demonstrate reactivity
	time.Sleep(500 * time.Millisecond)
	data1.Items[0].Message = "Baz"

	time.Sleep(500 * time.Millisecond)
	data1.Items[1].Message = "Qux"

	// Add a slice element.
	time.Sleep(500 * time.Millisecond)
	// I'm surprised that this actually works, but it does appear to.  I tried
	// appending 1000 items, and they're all reactive in the usual way.
	// Whatever magic GopherJS and Vue do appears to work together.
	data1.Items = append(data1.Items, NewItem("Quux"))

	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		data1.Items = append(data1.Items, NewItem(strconv.Itoa(i)))
	}
	time.Sleep(500 * time.Millisecond)
	randomSlice := data1.Items[3:5]
	randomSlice[1].Message = "I am randomSlice[1]"

	time.Sleep(500 * time.Millisecond)
	// Here's the way I would have expected that you'd have to do the above
	// appending -- slightly longer and more vulnerable to typos (since the
	// field name is just a string):
	hvue.Push(data1.Object.Get("items"), NewItem("Quux"))
}

func NewItem(m string) *ItemT {
	i := &ItemT{Object: O()}
	i.Message = m
	return i
}

package main

// From https://vuejs.org/v2/guide/computed.html#Computed-vs-Watched-Property
// and https://vuejs.org/v2/guide/computed.html#Computed-Setter.

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

type Data struct {
	*js.Object
	FirstName string `js:"firstName"`
	LastName  string `js:"lastName"`
	// FullName is computed

	VM *hvue.VM // Set by NewVM if you use DataS
}

func main() {
	d := &Data{Object: hvue.NewObject()}
	d.FirstName = "Foo"
	d.LastName = "Bar"

	vm := hvue.NewVM(
		hvue.El("#demo"),
		hvue.DataS(d),
		hvue.ComputedWithGetSet(
			"fullName",
			func(*hvue.VM) interface{} {
				return d.FirstName + " " + d.LastName
			},
			func(_ *hvue.VM, newValue *js.Object) {
				names := strings.Fields(newValue.String())
				d.FirstName = names[0]
				d.LastName = names[len(names)-1]
			}))

	time.Sleep(time.Second)
	vm.Set("fullName", "First Middle Last")
	// Note that FirstName & LastName are changed, too, and that "Middle"
	// is effectively ignored: the full value is not stored, but computed
	// from FirstName & LastName.
}

func (d *Data) FullName() string {
	return d.VM.Get("fullName").String()
}

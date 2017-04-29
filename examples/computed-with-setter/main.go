package main

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

type Data struct {
	*js.Object
	FirstName string `js:"firstName"`
	LastName  string `js:"lastName"`
	// FullName is computed
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
	go func() {
		time.Sleep(time.Second)
		vm.Set("fullName", "Foo Bar Baz")
	}()
}

package main

// From https://vuejs.org/v2/guide/computed.html#Computed-vs-Watched-Property
// and https://vuejs.org/v2/guide/computed.html#Computed-Setter.

import (
	"strings"
	"time"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"

	"github.com/huckridgesw/hvue"
)

type Data struct {
	js.Value

	VM *hvue.VM // Set by NewVM if you use DataS
}

func (d *Data) FirstName() string       { return d.Get("firstName").String() }
func (d *Data) LastName() string        { return d.Get("lastName").String() }
func (d *Data) SetFirstName(new string) { d.Set("firstName", new) }
func (d *Data) SetLastName(new string)  { d.Set("lastName", new) }
func (d *Data) SetFullName(new string) {
	names := strings.Fields(new)
	d.SetFirstName(names[0])
	d.SetLastName(names[len(names)-1])
}

func main() {
	d := &Data{Value: hvue.NewObject()}
	d.SetFullName("Foo Bar")

	hvue.NewVM(
		hvue.El("#demo"),
		hvue.DataS(d, d.Value),
		hvue.ComputedWithGetSet(
			"fullName",
			func(*hvue.VM) interface{} {
				return d.FirstName() + " " + d.LastName()
			},
			func(_ *hvue.VM, newValue js.Value) {
				names := strings.Fields(newValue.String())
				d.SetFirstName(names[0])
				d.SetLastName(names[len(names)-1])
			}),
		hvue.MethodsOf(&Data{}),
	)

	time.Sleep(time.Second)
	d.SetFullName("First Middle Last")
	// Note that FirstName & LastName are changed, too, and that "Middle"
	// is effectively ignored: the full value is not stored, but computed
	// from FirstName & LastName.

	select {}
}

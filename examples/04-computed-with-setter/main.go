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

func (d *Data) FirstName() string { return d.Get("firstName").String() }
func (d *Data) LastName() string  { return d.Get("lastName").String() }
func (d *Data) FullName() string  { return d.Get("fullName").String() }

func (d *Data) SetFirstName(new string) { d.Set("firstName", new) }
func (d *Data) SetLastName(new string)  { d.Set("lastName", new) }
func (d *Data) CalcFullName()           { d.Set("fullName", d.FirstName()+" "+d.LastName()) }
func (d *Data) SetFullName(new string) {
	names := strings.Fields(new)
	d.SetFirstName(names[0])
	d.SetLastName(names[len(names)-1])
	d.CalcFullName()
}

func main() {
	d := &Data{Value: hvue.NewObject()}
	d.SetFullName("Foo Bar")

	hvue.NewVM(
		hvue.El("#demo"),
		hvue.DataS(d, d.Value),
		// ComputedWithGetSet not implemented
		// hvue.ComputedWithGetSet(
		// 	"fullName",
		// 	func(*hvue.VM) interface{} {
		// 		return d.FirstName + " " + d.LastName
		// 	},
		// 	func(_ *hvue.VM, newValue js.Value) {
		// 		names := strings.Fields(newValue.String())
		// 		d.FirstName = names[0]
		// 		d.LastName = names[len(names)-1]
		// 	})
		hvue.MethodsOf(&Data{}),
		hvue.Watch("firstName", func(*hvue.VM) { d.CalcFullName() }),
		hvue.Watch("lastName", func(*hvue.VM) { d.CalcFullName() }),
	)

	time.Sleep(time.Second)
	d.SetFullName("First Middle Last")
	// Note that FirstName & LastName are changed, too, and that "Middle"
	// is effectively ignored: the full value is not stored, but computed
	// from FirstName & LastName.

	select {}
}

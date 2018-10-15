package main

// From https://vuejs.org/v2/guide/computed.html#Basic-Example

import (
	"strings"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"

	"github.com/huckridgesw/hvue"
)

type Data struct {
	js.Value
}

func (d *Data) Message() string       { return d.Get("message").String() }
func (d *Data) SetMessage(new string) { d.Set("message", new) }

func main() {
	d := &Data{Value: hvue.NewObject()}
	d.SetMessage("Hello")

	app := hvue.NewVM(
		hvue.El("#example"),
		hvue.DataS(d, d.Value),
		hvue.Computed(
			"reversedMessage",
			func(*hvue.VM) interface{} {
				return reverse(d.Message())
			}),
	)
	js.Global().Set("app", app.Value)
	// In the JS console, try app.message = "some other string"
	// Browser should change to Computed reversed message: "gnirts rehto emos".

	select {}
}

func reverse(s string) string {
	runes := strings.Split(s, "")
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return strings.Join(runes, "")
}

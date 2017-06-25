package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

type Data struct {
	*js.Object
	Message string `js:"message"`
}

func main() {
	d := &Data{Object: hvue.NewObject()}
	d.Message = "Hello"

	hvue.NewVM(
		hvue.El("#example"),
		hvue.DataS(d),
		hvue.Computed(
			"reversedMessage",
			func(vm *hvue.VM) interface{} {
				return reverse(d.Message)
			}))
}

func reverse(s string) string {
	runes := strings.Split(s, "")
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return strings.Join(runes, "")
}

package main

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

type Data struct {
	*js.Object
	Message string `js:"message"`
}

func main() {
	hvue.NewVM(
		hvue.El("#app-5"),
		hvue.DataS(hvue.NewT(&Data{Message: "Hello, Vue!"})),
		hvue.MethodsOf(&Data{}))
}

func (d *Data) ReverseMessage() {
	d.Message = reverse(d.Message)
}

func reverse(s string) string {
	runes := strings.Split(s, "")
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return strings.Join(runes, "")
}

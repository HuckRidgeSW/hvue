package main

import (
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
		hvue.DataS(NewData("Hello, Vue!")),
		hvue.MethodsOf(&Data{}))
}

func NewData(message string) *Data {
	d := &Data{Object: js.Global.Get("Object").New()}
	d.Message = message
	return d
}

func (d *Data) ReverseMessage(event *js.Object) {
	// event ignored
	d.Message = reverse(d.Message)
}

func reverse(s string) string {
	b := []byte(s)
	newS := make([]byte, len(s))
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		newS[i], newS[j] = b[j], b[i]
	}
	return string(newS)
}

package hvue

import "github.com/gopherjs/gopherjs/js"

func NewObject() *js.Object {
	return js.Global.Get("Object").New()
}

func NewArray() *js.Object {
	return js.Global.Get("Array").New()
}

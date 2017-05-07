package hvue

import "github.com/gopherjs/gopherjs/js"

func NewObject() *js.Object {
	return js.Global.Get("Object").New()
}

func NewArray() *js.Object {
	return js.Global.Get("Array").New()
}

// Append in place to the end of an array
func Push(o *js.Object, any interface{}) (newLength int) {
	return o.Call("push", any).Int()
}

// Vue.set
func Set(o, key, value interface{}) interface{} {
	js.Global.Get("Vue").Call("set", o, key, value)
	return value
}

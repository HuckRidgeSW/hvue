package hvue

import (
	"github.com/gopherjs/gopherwasm/js"
	// "github.com/gopherjs/gopherjs/js"
)

// NewObject is a utility function for creating a new JavaScript Object of
// type js.Value.
func NewObject() js.Value {
	return js.Global().Get("Object").New()
}

// NewArray is a utility function for creating a new JS array.
func NewArray() js.Value {
	return js.Global().Get("Array").New()
}

// Push appends any to the end of o, in place.
func Push(o js.Value, any interface{}) (newLength int) {
	return o.Call("push", any).Int()
}

// Set is a wrapper for js{Vue.set}
func Set(o, key, value interface{}) interface{} {
	js.Global().Get("Vue").Call("set", o, key, value)
	return value
}

func jsCallWithVM(f func(*VM) interface{}) js.Value {
	return NewCallback(
		func(this js.Value, args []js.Value) interface{} {
			vm := &VM{Value: this}
			return f(vm)
		})
}

func NewCallback(f func(this js.Value, args []js.Value) interface{}) js.Value {
	return js.Global().Call("wasm_callback",
		js.NewCallback(func(args []js.Value) {
			f(args[0], args[1:])
		}))
}

package hvue

import (
	// "github.com/gopherjs/gopherwasm/js"
	"fmt"
	"math"
	"syscall/js"
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

func jsCallWithVM(f func(*VM) interface{}) js.Callback {
	return js.NewCallback(
		func(this js.Value, args []js.Value) interface{} {
			vm := &VM{Value: this}
			return f(vm)
		})
}

func Log(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}

func GetDeep(o js.Value, fields ...string) (js.Value, error) {
	for _, field := range fields {
		new := o.Get(field)
		if new == js.Undefined() {
			return js.Value{}, fmt.Errorf("GetDeep: Empty field %s", field)
		}
		o = new
	}
	return o, nil
}

// In JavaScript, a truthy value is a value that is considered true when
// encountered in a Boolean context. All values are truthy unless they are
// defined as falsy (i.e., except for false, 0, "", null, undefined, and NaN).
func Falsy(o js.Value) bool {
	return !Truthy(o)
}

func Truthy(o js.Value) bool {
	switch o.Type() {
	case js.TypeUndefined, js.TypeNull:
		return false
	case js.TypeBoolean:
		return o.Bool()
	case js.TypeNumber:
		if math.IsNaN(o.Float()) {
			return false
		}
		return o.Float() != 0
	case js.TypeString:
		return o.String() != ""
	case js.TypeSymbol, js.TypeObject, js.TypeFunction:
		return true
	}
	return true
}

package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

// NewObject is a utility function for creating a new *js.Object.
func NewObject() *js.Object {
	return js.Global.Get("Object").New()
}

// NewArray is a utility function for creating a new JS array.
func NewArray() *js.Object {
	return js.Global.Get("Array").New()
}

// Push appends any to the end of o, in place.
func Push(o *js.Object, any interface{}) (newLength int) {
	return o.Call("push", any).Int()
}

// Set is a wrapper for js{Vue.set}
func Set(o, key, value interface{}) interface{} {
	js.Global.Get("Vue").Call("set", o, key, value)
	return value
}

// NewT is an attempt at making GopherJS struct initialization easier and more
// Go-like.  The intent is that instead of saying
//
//    f := T{Object: NewObject()}
//    f.slot1 = foo
//    f.slot2 = bar
//    // etc
//
// you can say
//
//    f := NewT{&T{slot1: foo, slot2: bar}}
//
// t should be a pointer, as in the examples above.
//
// Warning: Experimental.  Use at your own risk, and write a unit test.
// Doesn't work with complex types.  Not sure exactly what that means yet, but
// it means at least that it doesn't work with fields of type []string, or
// probably anything that's not a "basic" type such as int, float, string,
// etc.
func NewT(t interface{}) interface{} {
	io := js.InternalObject(t)
	valueOfT := reflect.ValueOf(t).Elem()

	// If the first field (assumed to be the *js.Object field) is set, just
	// return t unchanged.  Does no other error checking.  Should really check
	// for non-js fields and panic if it finds them.
	f0Name := valueOfT.Type().Field(0).Name
	if io.Get(f0Name) != nil {
		return t
	}

	if !valueOfT.Field(0).CanSet() {
		// reflect's Set method won't set unexported fields
		panic("The *js.Object field must be exported")
	}

	typ := valueOfT.Type()
	obj := o()

	for field := 1; field < typ.NumField(); field++ {
		if jsName, ok := typ.Field(field).Tag.Lookup("js"); ok {
			goName := typ.Field(field).Name
			obj.Set(jsName, io.Get(goName))
		}
	}

	valueOfT.Field(0).Set(reflect.ValueOf(obj))
	return t
}

func jsCallWithVM(f func(*VM) interface{}) *js.Object {
	return js.MakeFunc(
		func(this *js.Object, args []*js.Object) interface{} {
			vm := &VM{Object: this}
			return f(vm)
		})
}

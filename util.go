package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

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

// Warning: Doesn't work with complex types.  Not sure exactly what that means
// yet, but it means at least that it doesn't work with fields of type
// []string, or probably anything that's not a "basic" type such as int,
// float, string, etc.
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

package hvue

import (
	"reflect"
	"unsafe"

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
//    f := &T{Object: NewObject()}
//    f.Slot1 = foo
//    f.Slot2 = bar
//    // etc
//
// you can say
//
//    f := NewT(&T{Slot1: foo, Slot2: bar}).(*T)
//
// t should be a pointer, as in the examples above.  Only exported fields are
// set.
//
// Warning: Experimental.  Use at your own risk, and write a unit test.
// Doesn't work with complex types.  Not sure exactly what that means yet, but
// it means at least that it doesn't work with fields of type []string, or
// probably anything that's not a "basic" type such as int, float, string, etc.
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

	for field, nFields := 1, typ.NumField(); field < nFields; field++ {
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

// Cloak encapsulates a Go value within a JavaScript object. None of the fields
// or methods of the value will be exposed; it is therefore not intended that this
// *Object be used by Javascript code. Instead this function exists as a convenience
// mechanism for carrying state from Go to JavaScript and back again.
//
// Credit to Paul Jolly; see https://github.com/gopherjs/gopherjs/issues/704#issuecomment-332109410
func Cloak(i interface{}) *js.Object {
	return js.InternalObject(i)
}

// Uncloak is the inverse of Cloak.
func Uncloak(o *js.Object) interface{} {
	return interface{}(unsafe.Pointer(o.Unsafe()))
}

// 	v := S{Name: "Rob Pike"}
//
// 	// ...use js.InternalObject() on a _pointer_ to that value.
// 	//
// 	// Here we set a global variable in Javascript world to
// 	// the result.
// 	js.Global.Set("banana", js.InternalObject(&v))
//
// 	// ....
//
// 	// To then use that value again, get the *js.Object value...
// 	vjo := js.Global.Get("banana")
//
// 	// ... then use .Unsafe() + unsafe.Pointer
// 	vp := (*S)(unsafe.Pointer(vjo.Unsafe()))
//
// 	// Verify that we have exactly the same object
// 	println(vp == &v, *vp == v)
// }

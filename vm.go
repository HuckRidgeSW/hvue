package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type VM struct {
	*js.Object
}

type Config struct {
	*js.Object
	El      string     `js:"el"`
	Data    *js.Object `js:"data"`
	Methods *js.Object `js:"methods"`
	// Template string     `js:"template"`
}

type option func(*Config)

// NewVM returns a new vm, analogous to Javascript `new Vue(...)`.  See
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis and
// https://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html
// for discussions of the options.
func NewVM(opts ...option) *VM {
	c := &Config{Object: js.Global.Get("Object").New()}
	c.Option(opts...)
	return &VM{Object: js.Global.Get("Vue").New(c)}
}

// El sets the vm's el slot.
func El(selector string) option {
	return func(c *Config) {
		c.El = selector
	}
}

// Data sets a single data field.  Data can be called multiple times for the
// same vm.
func Data(name string, value interface{}) option {
	return func(c *Config) {
		if c.Data == js.Undefined {
			c.Data = js.Global.Get("Object").New()
		}
		c.Data.Set(name, value)
	}
}

// DataS sets the struct `value` as the entire contents of the vm's data
// field.
func DataS(value interface{}) option {
	return func(c *Config) {
		if c.Data != js.Undefined {
			panic("Cannot use hvue.Data and hvue.DataS together")
			c.Data = js.Global.Get("Object").New()
		}
		c.Object.Set("data", value)
	}
}

// MethodsOf sets up vm.methods with the exported methods of the type that t
// is an instance of.  Call it like MethodsOf(&SomeType{}).  SomeType must be
// a pure Javascript object, with no Go fields.  That is, all slots just have
// `js:"..."` tags.
func MethodsOf(t interface{}) option {
	return func(c *Config) {
		if c.Methods == js.Undefined {
			c.Methods = js.Global.Get("Object").New()
		}
		// Get the type of t
		typ := reflect.TypeOf(t)

		if typ.Kind() != reflect.Ptr {
			panic("Item passed to MethodsOf must be a pointer")
		}

		// Create a new receiver.  "Same" receiver used for all methods, with
		// its Object slot set differently(?) each time.  typ is a pointer type
		// so you have to get the type of the thing it points to with Elem() and
		// create a new one of those.
		receiver := reflect.New(typ.Elem())

		// Loop through all methods of the type
		for i := 0; i < typ.NumMethod(); i++ {
			// Get the i'th method's reflect.Method
			m := typ.Method(i)

			c.Methods.Set(m.Name,
				func(event *js.Object) {
					// Set the receiver's Object slot to c.Data.  receiver is a
					// pointer so you have to dereference it with Elem().
					receiver.Elem().Field(0).Set(reflect.ValueOf(c.Data))

					m.Func.Call([]reflect.Value{
						receiver,
						reflect.ValueOf(event)})
				})
		}
	}
}

// func Template(t string) option {
// 	return func(c *Config) {
// 		c.Template = t
// 	}
// }

// Option sets the options specified.
func (f *Config) Option(opts ...option) {
	for _, opt := range opts {
		opt(f)
	}
}

package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

var o = func() *js.Object { return js.Global.Get("Object").New() }

type VM struct {
	*js.Object
}

var jsOType = reflect.TypeOf(o())
var vmType = reflect.TypeOf(&VM{})

// NewVM returns a new vm, analogous to Javascript `new Vue(...)`.  See
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis and
// https://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html
// for discussions of how the options work, and also see the examples tree.
//
// If you use a data object (via DataS) and it has a VM field, it's set to
// this new VM.  TODO: Verify that the VM field is of type *hvue.VM.
func NewVM(opts ...option) *VM {
	c := &Config{Object: NewObject()}
	c.Option(opts...)
	vm := &VM{Object: js.Global.Get("Vue").New(c)}
	if c.dataValue.IsValid() {
		if vmField := c.dataValue.FieldByName("VM"); vmField.IsValid() {
			vmField.Set(reflect.ValueOf(vm))
		}
	}
	return vm
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
			c.Data = NewObject()
		}
		c.Data.Set(name, value)
	}
}

// DataS sets the struct `value` as the entire contents of the vm's data
// field.  `value` should be a pointer to the struct.  If the object has a VM
// field, NewVM sets it to the new VM object.
func DataS(value interface{}) option {
	return func(c *Config) {
		if c.Data != js.Undefined {
			panic("Cannot use hvue.DataS together with any other Data* options")
		}
		c.Object.Set("data", value)
		c.dataValue = reflect.ValueOf(value).Elem()
	}
}

func DataFunc(f func(*VM) interface{}) option {
	return func(c *Config) {
		if c.Data != js.Undefined {
			panic("Cannot use hvue.DataFunc together with any other Data* options")
		}

		c.Object.Set("data", js.MakeFunc(
			func(this *js.Object, jsArgs []*js.Object) interface{} {
				vm := &VM{Object: this}
				return f(vm)
			}))
	}
}

// MethodsOf sets up vm.methods with the exported methods of the type that t
// is an instance of.  Call it like MethodsOf(&SomeType{}).  SomeType must be
// a pure Javascript object, with no Go fields.  That is, all slots just have
// `js:"..."` tags.
//
// If a method wants a pointer to its vm, use a *VM as the first argument.
func MethodsOf(t interface{}) option {
	return func(c *Config) {
		if c.Methods == js.Undefined {
			c.Methods = NewObject()
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

			// Pre-compute some stuff that'd be the same for all calls of this
			// method.
			numIn := m.Type.NumIn()
			// If the 2nd arg (the *first* arg if you don't count the receiver)
			// expects a *VM, pass `this`.
			doVM := numIn > 1 && m.Type.In(1) == vmType

			c.Methods.Set(m.Name,
				js.MakeFunc(
					func(this *js.Object, jsArgs []*js.Object) interface{} {
						// Set the receiver's Object slot to c.Data.  receiver is a
						// pointer so you have to dereference it with Elem().
						receiver.Elem().Field(0).Set(reflect.ValueOf(c.Data))

						// Construct the arglist
						goArgs := make([]reflect.Value, numIn)
						goArgs[0] = receiver
						i := 1

						if doVM {
							vm := &VM{Object: this}
							goArgs[1] = reflect.ValueOf(vm)
							i++
						}

						for j := 0; j < len(jsArgs) && i < numIn; i, j = i+1, j+1 {
							switch m.Type.In(i).Kind() {
							case reflect.Ptr:
								inPtrType := m.Type.In(i)
								if inPtrType == jsOType {
									// A *js.Object
									goArgs[i] = reflect.ValueOf(jsArgs[j])
								} else {
									// Expects a pointer to a struct with first field
									// of type *js.Object.
									inType := inPtrType.Elem()
									inArg := reflect.New(inType)
									inArg.Elem().Field(0).Set(reflect.ValueOf(jsArgs[j]))
									goArgs[i] = inArg
								}
							case reflect.String:
								goArgs[i] = reflect.ValueOf(jsArgs[j].String())
							case reflect.Bool:
								goArgs[i] = reflect.ValueOf(jsArgs[j].Bool())
							case reflect.Float64:
								goArgs[i] = reflect.ValueOf(jsArgs[j].Float())
							case reflect.Int32, reflect.Int:
								goArgs[i] = reflect.ValueOf(jsArgs[j].Int())
							case reflect.Int64:
								goArgs[i] = reflect.ValueOf(jsArgs[j].Int64())
							default:
								panic("Unknown type in arglist for " +
									m.Name + ": " + m.Type.In(i).Kind().String())
							}
						}

						result := m.Func.Call(goArgs)

						// I don't think method results are ever actually used, but
						// I could be wrong.
						if len(result) >= 1 {
							return result[0].Interface()
						}
						return nil
					}))
		}
	}
}

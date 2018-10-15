package hvue

import (
	"reflect"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"
)

// VM wraps a js Vue object.
type VM struct {
	js.Value
}

func (vm *VM) Data() js.Value    { return vm.Get("$data") }
func (vm *VM) Props() js.Value   { return vm.Get("$props") }
func (vm *VM) El() js.Value      { return vm.Get("$el") }
func (vm *VM) Options() js.Value { return vm.Get("$options") }
func (vm *VM) Parent() js.Value  { return vm.Get("$parent") }
func (vm *VM) Root() js.Value    { return vm.Get("$root") }

// func (vm *VM) Children() []js.Value    { return vm.Get("$children") } // not sure about this one
func (vm *VM) Slots() js.Value       { return vm.Get("$slots") }
func (vm *VM) ScopedSlots() js.Value { return vm.Get("$scopedSlots") }
func (vm *VM) IsServer() bool        { return vm.Get("$isServer").Bool() }

// Note existence of fields with setter methods, which won't show up in $data.
func (vm *VM) Setters() js.Value { return vm.Get("hvue_setters") }

func (vm *VM) SetSetters(new js.Value) { vm.Value.Set("hvue_setters", new) }

var (
	jsOType     = reflect.TypeOf(NewObject())
	vmType      = reflect.TypeOf(&VM{})
	dataObjects = map[int]interface{}{}
	nextDataID  = 1
)

// NewVM returns a new vm, analogous to Javascript `new Vue(...)`.  See
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis and
// https://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html
// for discussions of how the options work, and also see the examples tree.
//
// If you use a data object (via DataS) and it has a VM field, it's set to
// this new VM.  TODO: Verify that the VM field is of type *hvue.VM.
func NewVM(opts ...ComponentOption) *VM {
	c := &Config{Value: NewObject()}
	c.SetSetters(NewObject())
	c.Option(opts...)
	cv := js.Global().Get("Vue").New(c.Value)
	vm := &VM{Value: cv}
	if c.dataValue.IsValid() {
		if vmField := c.dataValue.FieldByName("VM"); vmField.IsValid() {
			vmField.Set(reflect.ValueOf(vm))
		}
	}
	vm.SetSetters(c.Setters())
	return vm
}

// El sets the vm's el slot.
func El(selector string) ComponentOption {
	return func(c *Config) {
		c.SetEl(selector)
	}
}

// Data sets a single data field.  Data can be called multiple times for the
// same vm.
//
// Note that you can't use MethodsOf with this function.
func Data(name string, value interface{}) ComponentOption {
	return func(c *Config) {
		if c.Data() == js.Undefined() {
			c.SetData(NewObject())
		}
		c.Data().Set(name, value)
	}
}

// DataS sets the object `goValue` as the entire contents of the vm's data
// field.  If the object has a VM field, NewVM sets it to the new VM object.
func DataS(goValue interface{}, jsValue js.Value) ComponentOption {
	return func(c *Config) {
		if c.Data() != js.Undefined() {
			panic("Cannot use hvue.DataS together with any other Data* options")
		}
		c.SetData(jsValue)
		c.dataValue = reflect.ValueOf(goValue).Elem()
		storeDataID(jsValue, goValue, c)
	}
}

// DataFunc defines a function that returns a new data object.  You have to
// use DataFunc with Components, not Data or DataS.
//
// If you use DataFunc and MethodsOf together, the type of the object returned
// by DataFunc should match the type of the object given to MethodsOf.
//
// Note that this function is called when the VM or component is created
// (https://vuejs.org/v2/api/#created), not when you call "NewVM".  This means
// that you can't, for example, get clever and try to use the same object here
// as with MethodsOf.  MethodsOf requires an object when you call NewVM to
// register the VM, long before the VM is actually created or bound; this is
// called every time a new VM or component is created.
func DataFunc(f DataFuncT) ComponentOption {
	return func(c *Config) {
		if c.Data() != js.Undefined() {
			panic("Cannot use hvue.DataFunc together with any other Data/DataS options")
		}
		c.SetDataFunc(f)
	}
}

// Store a data object ID in the data object, for later reference.
//
// This wouldn't work if the js.Value is sealed or not "plain" (like
// WebSocket).  But on the other hand, Vue won't work with non-plain or sealed
// objects, so it doesn't matter.
func storeDataID(jsValue js.Value, goValue interface{}, c *Config) {
	curID := nextDataID // small race condition here
	nextDataID++
	jsValue.Set("hvue_dataID", curID)

	// Store the Go data object, indexed by curID
	dataObjects[curID] = goValue

	// Schedule it to be deleted when the vm is deleted
	Destroyed(func(*VM) {
		delete(dataObjects, curID)
	})(c)

}

// Method adds a single function as a "method" on a vm.  It does not change
// the method set of the data object, if any.
func Method(name string, f interface{}) ComponentOption {
	return func(c *Config) {
		if c.Methods() == js.Undefined() {
			c.SetMethods(NewObject())
		}
		m := reflect.ValueOf(f)
		if m.Kind() != reflect.Func {
			panic("Method " + name + " is not a func")
		}

		c.Methods().Set(name,
			makeMethod(name, false, m.Type(), m))
	}
}

// MethodsOf sets up vm.methods with the exported methods of the type that t
// is an instance of.  Call it like MethodsOf(&SomeType{}).  SomeType must be
// a pure Javascript object, with no Go fields.  That is, all slots just have
// `js:"..."` tags.
//
// If a method wants a pointer to its vm, use a *VM as the first argument.
//
// You can't use MethodsOf with Data(), only with DataS or DataFunc().
func MethodsOf(t interface{}) ComponentOption {
	return func(c *Config) {
		if c.Methods() == js.Undefined() {
			c.SetMethods(NewObject())
		}
		typ := reflect.TypeOf(t)
		if typ.Kind() != reflect.Ptr ||
			typ.Elem().Kind() != reflect.Struct {
			panic("Item passed to MethodsOf must be a pointer to a struct")
		}

		// Loop through all methods of the type
		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			c.Methods().Set(m.Name,
				makeMethod(m.Name, true, m.Type, m.Func))
		}
	}
}

func makeMethod(name string, isMethod bool, mType reflect.Type, m reflect.Value) js.Callback {
	return js.NewCallback(
		func(this js.Value, jsArgs []js.Value) interface{} {
			// Construct the arglist
			numIn := mType.NumIn()
			goArgs := make([]reflect.Value, numIn)
			goArg := 0

			if isMethod {
				// Lookup the receiver in dataObjects, based on
				// $data.hvue_dataID
				dataID := this.Get("$data").Get("hvue_dataID").Int()
				if dataID == 0 {
					// FIXME: A better error here would be great, Mmmkay?
					panic("Unknown dataID for method " + name)
				}
				receiver, ok := dataObjects[dataID]
				if !ok {
					panic("Unknown dataID for method " + name)
				}

				goArgs[0] = reflect.ValueOf(receiver)
				goArg = 1
			}

			vmDone := false
			// We say || in the WHILE clause instead of && because there could be
			// Go args (like the receiver and a *VM arg) that wouldn't show up in
			// the JS arglist.
			for jsArg := 0; jsArg < len(jsArgs) || goArg < numIn; goArg, jsArg = goArg+1, jsArg+1 {
				if goArg >= numIn {
					break
				}

				switch mType.In(goArg) {
				case jsOType:
					// A js.Value
					goArgs[goArg] = reflect.ValueOf(jsArgs[jsArg])
				default:
					switch mType.In(goArg).Kind() {
					case reflect.Ptr:
						inPtrType := mType.In(goArg)
						switch inPtrType {
						case vmType:
							// A *VM
							if vmDone {
								panic("Only a single *hvue.VM arg expected per method: " + name)
							}
							goArgs[goArg] = reflect.ValueOf(&VM{Value: this})
							jsArg--
							vmDone = true
						default:
							// Expects a pointer to a struct with first field
							// of type js.Value.  Doesn't work yet with nested
							// structs.
							inType := inPtrType.Elem()
							inArg := reflect.New(inType)
							inArg.Elem().Field(0).Set(reflect.ValueOf(jsArgs[jsArg]))
							goArgs[goArg] = inArg
						}
					case reflect.String:
						goArgs[goArg] = reflect.ValueOf(jsArgs[jsArg].String())
					case reflect.Bool:
						goArgs[goArg] = reflect.ValueOf(jsArgs[jsArg].Bool())
					case reflect.Float64:
						goArgs[goArg] = reflect.ValueOf(jsArgs[jsArg].Float())
					case reflect.Int64, reflect.Int32, reflect.Int:
						goArgs[goArg] = reflect.ValueOf(jsArgs[jsArg].Int())
					default:
						panic("hvue.makeMethod: Unknown type in arglist for " +
							name + ": " + mType.In(goArg).Kind().String())
					}
				}
			}

			result := m.Call(goArgs)

			// I don't think method results are ever actually used, but
			// I could be wrong.
			if len(result) >= 1 {
				return result[0].Interface()
			}
			return nil
		})
}

func Watch(name string, f func(*VM)) ComponentOption {
	return func(c *Config) {
		if c.Watchers() == js.Undefined() {
			c.SetWatchers(NewObject())
		}

		c.Watchers().Set(
			name,
			jsCallWithVM(func(vm *VM) interface{} {
				f(vm)
				return nil
			}))
	}
}

// FIXME: A filter function needs to be able to return a value, which Go
// functions can't yet.  So comment this out for now.
// func Filter(name string, f func(vm *VM, value js.Value, args ...js.Value) interface{}) ComponentOption {
// 	return func(c *Config) {
// 		if c.Filters() == js.Undefined() {
// 			c.SetFilters(NewObject())
// 		}
//
// 		c.Filters().Set(name, js.NewCallback(
// 			func(args []js.Value) interface{} {
// 				vm := &VM{Value: args[0]}
// 				return f(vm, args[0], args[1:]...)
// 			}))
// 	}
// }

// Emit emits an event.  It wraps js{vm.$emit}:
// https://vuejs.org/v2/api/#vm-emit.
func (vm *VM) Emit(event string, args ...interface{}) {
	args = append([]interface{}{event}, args...)
	vm.Call("$emit", args...)
}

// Refs returns the ref for name.  vm.Refs("foo") compiles to
// js{vm.$refs.foo}.  It wraps vm.$refs: https://vuejs.org/v2/api/#vm-refs.
func (vm *VM) Refs(name string) js.Value {
	return vm.Get("$refs").Get(name)
}

// GetData returns the Go data object associated with a *VM.  You need to type
// assert its return value to data type you passed to DataS(), or returned
// from the function given to DataFunc().
func (vm *VM) GetData() interface{} {
	dataID := vm.Data().Get("hvue_dataID").Int()
	if dataID == 0 {
		// FIXME: A better error here would be great, Mmmkay?
		panic("Unknown dataID in GetData")
	}
	dataObj, ok := dataObjects[dataID]
	if !ok {
		panic("Unknown dataID in GetData")
	}
	return dataObj
}

// Set wraps vm.Value.Set(), but checks to make sure the given field is a
// valid slot in the VM's data object (including computed setters), and panics
// otherwise.  (If you don't want this check, then use vm.Value.Set()
// directly.)
func (vm *VM) Set(key string, value interface{}) {
	if vm.Data().Get(key) == js.Undefined() &&
		vm.Setters().Get(key) == js.Undefined() {
		panic("Unknown data slot set: " + key)
	}
	vm.Value.Set(key, value)
}

// Modeled on GopherJS's js.M, also a map[string]interface{}
type M map[string]interface{}

func Map2Obj(m M) js.Value {
	res := NewObject()
	for k, v := range m {
		if m, ok := v.(M); ok {
			res.Set(k, Map2Obj(m))
		} else {
			res.Set(k, v)
		}
	}
	return res
}

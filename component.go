package hvue

import (
	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"
)

// NewComponent defines a new Vue component.  It wraps js{Vue.component}:
// https://vuejs.org/v2/api/#Vue-component.
func NewComponent(name string, opts ...ComponentOption) {
	c := &Config{Value: NewObject()}
	c.SetSetters(NewObject())
	c.Option(opts...)

	if c.DataType == js.TypeUndefined {
		// wasm_new_data_func takes care of the hvue_dataID magic.
		// c.Set("data",
		// 	js.Global().Call("wasm_new_data_func",
		// 		NewObject(), // call wasm_new_data_func with a blank template
		// 		js.NewCallback(func([]js.Value) {}),
		// 	))

		c.Set("data",
			js.NewCallback(func(this js.Value, _ []js.Value) interface{} {
				newO := NewObject()
				dataID := this.Get("$parent").Get("$data").Get("hvue_dataID")
				if dataID != js.Undefined() {
					newO.Set("hvue_dataID", dataID.Int())
				}
				return newO
			}))
	} else if c.DataType != js.TypeFunction {
		panic("Cannot use Data() with NewComponent, must use DataFunc.  Component: " + name)
	}

	js.Global().Get("Vue").Call("component", name, c.Value)
}

// Component is used in NewVM to define a local component, within the scope of
// another instance/component.
// https://vuejs.org/v2/guide/components.html#Local-Registration
func Component(name string, opts ...ComponentOption) ComponentOption {
	return func(c *Config) {
		componentOption := &Config{Value: NewObject()}
		componentOption.Option(opts...)

		if c.Components() == js.Undefined() {
			c.SetComponents(NewObject())
		}

		c.Components().Set(name, componentOption.Value)
	}
}

// Props defines one or more simple prop slots.  For complex prop slots, use
// PropObj().  https://vuejs.org/v2/api/#props
func Props(props ...string) ComponentOption {
	return func(c *Config) {
		if c.Props() == js.Undefined() {
			c.SetProps(NewArray())
		}
		for i, prop := range props {
			c.Props().SetIndex(i, prop)
		}
	}
}

// PropObj defines a complex prop slot called `name`, configured with Types,
// Default, DefaultFunc, and Validator.
func PropObj(name string, opts ...PropOption) ComponentOption {
	return func(c *Config) {
		if c.Props() == js.Undefined() {
			c.SetProps(NewObject())
		}
		pO := &PropConfig{Value: NewObject()}
		pO.Option(opts...)
		c.Props().Set(name, pO.Value)
	}
}

// Template defines a template for a component.  It sets the js{template} slot
// of a js{Vue.component}'s configuration object.
func Template(template string) ComponentOption {
	return func(c *Config) {
		c.SetTemplate(template)
	}
}

// Types configures the allowed types for a prop.
// https://vuejs.org/v2/guide/components.html#Props.
func Types(types ...pOptionType) PropOption {
	return func(p *PropConfig) {
		arr := NewArray()
		for _, t := range types {
			var newVal js.Value
			switch t {
			case PString:
				newVal = js.Global().Get("String")
			case PNumber:
				newVal = js.Global().Get("Number")
			case PBoolean:
				newVal = js.Global().Get("Boolean")
			case PFunction:
				newVal = js.Global().Get("Function")
			case PObject:
				newVal = js.Global().Get("Object")
			case PArray:
				newVal = js.Global().Get("Array")
			}
			arr.Call("push", newVal)
		}
		p.SetType(arr)
	}
}

// Required specifies that the prop is required.
// https://vuejs.org/v2/guide/components.html#Props.
var Required PropOption = func(p *PropConfig) {
	p.SetRequired(true)
}

// Default gives the default for a prop.
// https://vuejs.org/v2/guide/components.html#Props
func Default(def interface{}) PropOption {
	return func(p *PropConfig) {
		p.SetDefault(def)
	}
}

// DefaultFunc sets a function that returns the default for a prop.
// https://vuejs.org/v2/guide/components.html#Props
//
// FIXME: Right now, can only pass an object (not a function).  The JS helper
// function copies it to a new object.  Later, need to be able to pass a
// function, which returns a new value.
func DefaultFunc(def js.Value) PropOption {
	return func(p *PropConfig) {
		p.SetDefault(js.Global().Call("wasm_return_copy", def))
	}
}

// Validator functions generate warnings in the JS console if using the
// vue.js development build.  They don't panic or otherwise crash your code,
// they just give warnings if the validation fails.
//
// FIXME: Currently does nothing, because in 1.11 Go functions called from JS
// can't return values.
func Validator(f func(vm *VM, value js.Value) interface{}) PropOption {
	return func(p *PropConfig) {
		return
		// p.SetValidator(NewCallback(
		// 	func(this js.Value, args []js.Value) interface{} {
		// 		vm := &VM{Value: this}
		// 		return f(vm, args[0])
		// 	})
	}
}

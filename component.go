package hvue

import "github.com/gopherjs/gopherjs/js"

func NewComponent(name string, opts ...ComponentOption) {
	c := &Config{Object: o()}
	c.Option(opts...)

	if c.Data == js.Undefined {
		c.Object.Set("data", jsCallWithVM(func(vm *VM) interface{} {
			obj := o()
			// Get the parent data object ID, if it exists
			dataID := vm.Get("$parent").Get("$data").Get("hvue_dataID")
			if dataID != js.Undefined {
				obj.Set("hvue_dataID", dataID)
			}
			return obj
		}))
	}

	js.Global.Get("Vue").Call("component", name, c.Object)
}

func Component(name string, data interface{}) ComponentOption {
	return func(c *Config) {
		if c.Components == js.Undefined {
			c.Components = o()
		}
		c.Components.Set(name, data)
	}
}

func Props(props ...string) ComponentOption {
	return func(c *Config) {
		if c.Props == js.Undefined {
			c.Props = NewArray()
		}
		for i, prop := range props {
			c.Props.SetIndex(i, prop)
		}
	}
}

func PropObj(prop string, opts ...PropOption) ComponentOption {
	return func(c *Config) {
		if c.Props == js.Undefined {
			c.Props = o()
		}
		pO := &PropConfig{Object: o()}
		pO.Option(opts...)
		c.Props.Set(prop, pO.Object)
	}
}

func Template(template string) ComponentOption {
	return func(c *Config) {
		c.Template = template
	}
}

func Types(types ...pOptionType) PropOption {
	return func(p *PropConfig) {
		arr := js.Global.Get("Array").New()
		for _, t := range types {
			var newVal *js.Object
			switch t {
			case PString:
				newVal = js.Global.Get("String")
			case PNumber:
				newVal = js.Global.Get("Number")
			case PBoolean:
				newVal = js.Global.Get("Boolean")
			case PFunction:
				newVal = js.Global.Get("Function")
			case PObject:
				newVal = js.Global.Get("Object")
			case PArray:
				newVal = js.Global.Get("Array")
			}
			arr.Call("push", newVal)
		}
		p.typ = arr
	}
}

var Required PropOption = func(p *PropConfig) {
	p.required = true
}

func Default(def interface{}) PropOption {
	return func(p *PropConfig) {
		p.def = def
	}
}

func DefaultFunc(def func(*VM) interface{}) PropOption {
	return func(p *PropConfig) {
		p.def = jsCallWithVM(def)
	}
}

// Validator functions generate warnings in the JS console if using the
// vue.js development build.
func Validator(f func(vm *VM, value *js.Object) interface{}) PropOption {
	return func(p *PropConfig) {
		p.validator = js.MakeFunc(
			func(this *js.Object, args []*js.Object) interface{} {
				vm := &VM{Object: this}
				return f(vm, args[0])
			})
	}
}

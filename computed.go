package hvue

import "github.com/gopherjs/gopherjs/js"

func Computed(name string, f func(vm *VM) interface{}) option {
	return func(c *Config) {
		if c.Computed == js.Undefined {
			c.Computed = NewObject()
		}
		c.Computed.Set(name, js.MakeFunc(
			func(this *js.Object, _ []*js.Object) interface{} {
				vm := &VM{Object: this}
				return f(vm)
			}))
	}
}

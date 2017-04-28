package hvue

import "github.com/gopherjs/gopherjs/js"

func BeforeCreate(f func(vm *VM)) option  { return makeLifecycleMethod("beforeCreate", f) }
func Created(f func(vm *VM)) option       { return makeLifecycleMethod("created", f) }
func BeforeMount(f func(vm *VM)) option   { return makeLifecycleMethod("beforeMount", f) }
func Mounted(f func(vm *VM)) option       { return makeLifecycleMethod("mounted", f) }
func BeforeUpdate(f func(vm *VM)) option  { return makeLifecycleMethod("beforeUpdate", f) }
func Updated(f func(vm *VM)) option       { return makeLifecycleMethod("updated", f) }
func BeforeDestroy(f func(vm *VM)) option { return makeLifecycleMethod("beforeDestroy", f) }
func Destroyed(f func(vm *VM)) option     { return makeLifecycleMethod("destroyed", f) }

func makeLifecycleMethod(name string, f func(vm *VM)) option {
	return func(c *Config) {
		c.Set(name,
			js.MakeFunc(
				func(this *js.Object, _ []*js.Object) interface{} {
					vm := &VM{Object: this}
					f(vm)
					return nil
				}))
	}
}

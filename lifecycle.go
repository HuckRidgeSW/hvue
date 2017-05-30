package hvue

func BeforeCreate(f func(vm *VM)) ComponentOption  { return makeLifecycleMethod("beforeCreate", f) }
func Created(f func(vm *VM)) ComponentOption       { return makeLifecycleMethod("created", f) }
func BeforeMount(f func(vm *VM)) ComponentOption   { return makeLifecycleMethod("beforeMount", f) }
func Mounted(f func(vm *VM)) ComponentOption       { return makeLifecycleMethod("mounted", f) }
func BeforeUpdate(f func(vm *VM)) ComponentOption  { return makeLifecycleMethod("beforeUpdate", f) }
func Updated(f func(vm *VM)) ComponentOption       { return makeLifecycleMethod("updated", f) }
func BeforeDestroy(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("beforeDestroy", f) }
func Destroyed(f func(vm *VM)) ComponentOption     { return makeLifecycleMethod("destroyed", f) }

func makeLifecycleMethod(name string, f func(vm *VM)) ComponentOption {
	return func(c *Config) {
		c.Set(name,
			jsCallWithVM(func(vm *VM) interface{} {
				f(vm)
				return nil
			}))
	}
}

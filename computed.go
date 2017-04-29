package hvue

import "github.com/gopherjs/gopherjs/js"

// Define name as a computed property.  Note that name *must not* be set in
// data for this to work.  It's probably best if it's not even a slot in the
// struct.  Only access it via vm.Get/Set.
func Computed(name string, f func(vm *VM) interface{}) option {
	return func(c *Config) {
		if c.Computed == js.Undefined {
			c.Computed = NewObject()
		}
		c.Computed.Set(name,
			js.MakeFunc(
				func(this *js.Object, _ []*js.Object) interface{} {
					vm := &VM{Object: this}
					return f(vm)
				}))
	}
}

// Define name as a computed property with explicit get & set.  Note that name
// *must not* be set in data for this to work.  It's probably best if it's not
// even a slot in the struct.  Only access it via vm.Get/Set.
func ComputedWithGetSet(name string, get func(vm *VM) interface{}, set func(vm *VM, newValue *js.Object)) option {
	return func(c *Config) {
		if c.Computed == js.Undefined {
			c.Computed = NewObject()
		}
		c.Computed.Set(name,
			js.M{
				"get": js.MakeFunc(
					func(this *js.Object, _ []*js.Object) interface{} {
						vm := &VM{Object: this}
						return get(vm)
					}),
				"set": js.MakeFunc(
					func(this *js.Object, args []*js.Object) interface{} {
						vm := &VM{Object: this}
						set(vm, args[0])
						return nil
					})})
	}
}

package hvue

import (
	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"
)

// Computed defines name as a computed property.  Note that name *must not* be
// set in data for this to work.
//
// See the 04-computed-with-setter example.
func Computed(name string, f func(vm *VM) interface{}) ComponentOption {
	return func(c *Config) {
		if c.Computed() == js.Undefined() {
			c.SetComputed(NewObject())
		}
		c.Computed().Set(name, jsCallWithVM(f))
	}
}

// ComputedWithGetSet defines name as a computed property with explicit get &
// set.  Note that name *must not* be set in data for this to work.
//
// See the 04-computed-with-setter example.
func ComputedWithGetSet(name string, get func(vm *VM) interface{}, set func(vm *VM, newValue js.Value)) ComponentOption {
	return func(c *Config) {
		if c.Computed() == js.Undefined() {
			c.SetComputed(NewObject())
		}
		c.Computed().Set(name,
			map[string]interface{}{
				"get": jsCallWithVM(get),
				"set": js.NewCallback(
					func(this js.Value, args []js.Value) interface{} {
						vm := &VM{Value: this}
						set(vm, args[0])
						return nil
					})})
		c.Setters().Set(name, true)
	}
}

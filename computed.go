package hvue

import (
	"github.com/gopherjs/gopherwasm/js"
	// "github.com/gopherjs/gopherjs/js"
)

// Computed defines name as a computed property.  Note that name *must not* be
// set in data for this to work.  It's probably best if it's not even a slot
// in the struct.  Only access it via vm.Get.  You could also create an
// accessor; see the 04-computed-with-setter example.
func Computed(name string, f func(vm *VM) interface{}) ComponentOption {
	return func(c *Config) {
		if c.Computed() == js.Undefined() {
			c.SetComputed(NewObject())
		}
		c.Computed().Set(name, jsCallWithVM(f))
	}
}

// ComputedWithGetSet defines name as a computed property with explicit get &
// set.  Note that name *must not* be set in data for this to work.  It's
// probably best if it's not even a slot in the struct.  Only access it via
// vm.Get/Set.  You could create an accessor; see the 04-computed-with-setter
// example.
func ComputedWithGetSet(name string, get func(vm *VM) interface{}, set func(vm *VM, newValue js.Value)) ComponentOption {
	return func(c *Config) {
		if c.Computed() == js.Undefined() {
			c.SetComputed(NewObject())
		}
		c.Computed().Set(name,
			map[string]interface{}{
				"get": jsCallWithVM(get),
				"set": NewCallback(
					func(this js.Value, args []js.Value) interface{} {
						vm := &VM{Value: this}
						set(vm, args[0])
						return nil
					})})
		c.Setters().Set(name, true)
	}
}

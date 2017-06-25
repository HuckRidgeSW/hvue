package main

import (
	"time"

	"github.com/huckridge/hvue"
)

func main() {
	vm := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"),
		hvue.BeforeCreate(func(vm *hvue.VM) { println("BeforeCreate, vm:", vm.Object) }),
		hvue.Created(func(vm *hvue.VM) { println("Created") }),
		hvue.BeforeMount(func(vm *hvue.VM) { println("BeforeMount") }),
		hvue.Mounted(func(vm *hvue.VM) { println("Mounted") }),
		hvue.BeforeUpdate(func(vm *hvue.VM) { println("BeforeUpdate") }),
		hvue.Updated(func(vm *hvue.VM) { println("Updated") }),
		hvue.BeforeDestroy(func(vm *hvue.VM) { println("BeforeDestroy") }),
		hvue.Destroyed(func(vm *hvue.VM) { println("Destroyed") }))

	// Trigger the BeforeUpdate/Updated hooks
	time.Sleep(time.Second)
	vm.Set("message", "new data")

	// Trigger the BeforeDestroy/Destroyed hooks
	time.Sleep(time.Second)
	vm.Call("$destroy")

	// In the JS console, check for logs from the lifecycle.
}

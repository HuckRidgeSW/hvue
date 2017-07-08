package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridge/hvue"
)

func main() {
	vm := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"),
		hvue.Data("show", 1),
		hvue.BeforeCreate(func(vm *hvue.VM) { println("BeforeCreate, vm:", vm.Object) }),
		hvue.Created(func(vm *hvue.VM) { println("Created") }),
		hvue.BeforeMount(func(vm *hvue.VM) { println("BeforeMount") }),
		hvue.Mounted(func(vm *hvue.VM) { println("Mounted") }),
		hvue.BeforeUpdate(func(vm *hvue.VM) { println("BeforeUpdate") }),
		hvue.Updated(func(vm *hvue.VM) { println("Updated") }),
		hvue.BeforeDestroy(func(vm *hvue.VM) { println("BeforeDestroy") }),
		hvue.Destroyed(func(vm *hvue.VM) { println("Destroyed") }),
		hvue.Component("show1",
			hvue.Template("<div>custom component show1; show == 1</div>"),
			hvue.Activated(func(vm *hvue.VM) { println("Show1 activated") }),
			hvue.Deactivated(func(vm *hvue.VM) { println("Show1 deactivated") }),
		),
		hvue.Component("show2",
			hvue.Template("<div>custom component show2; show == 2</div>"),
			hvue.Activated(func(vm *hvue.VM) { println("Show2 activated") }),
			hvue.Deactivated(func(vm *hvue.VM) { println("Show2 deactivated") }),
		),
	)
	js.Global.Set("vm", vm)

	// Trigger the BeforeUpdate/Updated hooks
	time.Sleep(time.Second)
	vm.Set("message", "trigger the beforeUpdate/updated hooks")

	// Trigger the activated/deactivated hooks, which only run in components
	// inside a keep-alive.
	time.Sleep(time.Second)
	vm.Set("show", 2)
	time.Sleep(time.Second)
	vm.Set("show", 1)

	// Trigger the BeforeDestroy/Destroyed hooks
	time.Sleep(time.Second)
	vm.Call("$destroy")

	// In the JS console, check for logs from the lifecycle.
}

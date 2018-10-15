package main

// This file demos Vue lifecycle hooks
// (https://vuejs.org/v2/guide/instance.html#Instance-Lifecycle-Hooks) but does
// not correspond to any specific example on that page.
//
// See https://vuejs.org/v2/api/#Options-Lifecycle-Hooks for docs on all hooks.

import (
	"time"

	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"

	"github.com/huckridgesw/hvue"
)

func main() {
	vm := hvue.NewVM(
		hvue.El("#app"),
		hvue.Data("message", "Hello, Vue!"),
		hvue.Data("show", 1),
		hvue.BeforeCreate(func(vm *hvue.VM) { js.Global().Get("console").Call("log", "BeforeCreate, vm:", vm.Value) }),
		hvue.Created(func(*hvue.VM) { println("Created") }),
		hvue.BeforeMount(func(*hvue.VM) { println("BeforeMount") }),
		hvue.Mounted(func(*hvue.VM) { println("Mounted") }),
		hvue.BeforeUpdate(func(*hvue.VM) { println("BeforeUpdate") }),
		hvue.Updated(func(*hvue.VM) { println("Updated") }),
		hvue.BeforeDestroy(func(*hvue.VM) { println("BeforeDestroy") }),
		hvue.Destroyed(func(*hvue.VM) { println("Destroyed") }),
		hvue.Component("show1",
			hvue.Template("<div>custom component show1; show == 1</div>"),
			hvue.Activated(func(*hvue.VM) { println("Show1 activated") }),
			hvue.Deactivated(func(*hvue.VM) { println("Show1 deactivated") }),
		),
		hvue.Component("show2",
			hvue.Template("<div>custom component show2; show == 2</div>"),
			hvue.Activated(func(*hvue.VM) { println("Show2 activated") }),
			hvue.Deactivated(func(*hvue.VM) { println("Show2 deactivated") }),
		),
	)
	js.Global().Set("vm", vm.Value)

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

	select {}
}

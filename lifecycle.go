package hvue

// BeforeCreate lets you define a hook for the beforeCreate lifecycle action.
// "Called synchronously after the instance has just been initialized, before
// data observation and event/watcher setup."
// https://vuejs.org/v2/api/#beforeCreate
func BeforeCreate(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("beforeCreate", f) }

// Created lets you define a hook for the created lifecycle action.  "Called
// synchronously after the instance is created. At this stage, the instance
// has finished processing the options which means the following have been set
// up: data observation, computed properties, methods, watch/event callbacks.
// However, the mounting phase has not been started, and the $el property will
// not be available yet." https://vuejs.org/v2/api/#created
func Created(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("created", f) }

// BeforeMount lets you define a hook for the beforeMount lifecycle action.
// "Called right before the mounting begins: the render function is about to
// be called for the first time."  https://vuejs.org/v2/api/#beforeMount
func BeforeMount(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("beforeMount", f) }

// Mounted lets you define a hook for the mounted lifecycle action.  "Called
// after the instance has just been mounted where el is replaced by the newly
// created vm.$el. If the root instance is mounted to an in-document element,
// vm.$el will also be in-document when mounted is called."
// https://vuejs.org/v2/api/#mounted
func Mounted(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("mounted", f) }

// BeforeUpdate lets you define a hook for the beforeUpdate lifecycle action.
// "Called when the data changes, before the virtual DOM is re-rendered and
// patched.
//
// You can perform further state changes in this hook and they will not
// trigger additional re-renders."
// https://vuejs.org/v2/api/#beforeUpdate
func BeforeUpdate(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("beforeUpdate", f) }

// Updated lets you define a hook for the updated lifecycle action.  "Called
// after a data change causes the virtual DOM to be re-rendered and patched.
//
// The component’s DOM will have been updated when this hook is called, so you
// can perform DOM-dependent operations here. However, in most cases you
// should avoid changing state inside the hook. To react to state changes,
// it’s usually better to use a computed property or watcher instead."
// https://vuejs.org/v2/api/#updated
func Updated(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("updated", f) }

// Activated lets you define a hook for the activated lifecycle action.  Only
// runs in Vue-defined components (e.g. not regular DIVs) inside a <keep-alive>.
// "Called when a kept-alive component is activated."
// https://vuejs.org/v2/api/#activated and https://vuejs.org/v2/api/#keep-alive
func Activated(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("activated", f) }

// Deactivated lets you define a hook for the deactivated lifecycle action.
// Only runs in Vue-defined components (e.g. not regular DIVs) inside a
// <keep-alive>.
// "Called when a kept-alive component is deactivated."
// https://vuejs.org/v2/api/#deactivated and https://vuejs.org/v2/api/#keep-alive
func Deactivated(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("deactivated", f) }

// BeforeDestroy lets you define a hook for the beforeDestroy lifecycle
// action.  "Called right before a Vue instance is destroyed. At this stage
// the instance is still fully functional."
// https://vuejs.org/v2/api/#beforeDestroy
func BeforeDestroy(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("beforeDestroy", f) }

// Destroyed lets you define a hook for the destroyed lifecycle action.
// "Called after a Vue instance has been destroyed. When this hook is called,
// all directives of the Vue instance have been unbound, all event listeners
// have been removed, and all child Vue instances have also been destroyed."
// https://vuejs.org/v2/api/#destroyed
func Destroyed(f func(vm *VM)) ComponentOption { return makeLifecycleMethod("destroyed", f) }

func makeLifecycleMethod(name string, f func(vm *VM)) ComponentOption {
	return func(c *Config) {
		c.Set(name,
			jsCallWithVM(func(vm *VM) interface{} {
				f(vm)
				return nil
			}))
	}
}

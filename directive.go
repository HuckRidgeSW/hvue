package hvue

import (
	// "github.com/gopherjs/gopherwasm/js"
	"syscall/js"
)

// Directive wraps a js{Vue.directive} object.
// https://vuejs.org/v2/api/#Vue-directive.
type Directive struct {
	js.Value
}

// DirectiveBinding wraps the js{binding} slot of the directive hook argument.
// https://vuejs.org/v2/guide/custom-directive.html#Directive-Hook-Arguments
type DirectiveBinding struct {
	// This js.Value slot has its own slot called "value", so its accessor
	// (below) is called Value(), so the slot name can't also be Value, so call
	// it Val.  Which could actually be surprising if you use
	// DirectiveBinding.Value, because it'll *compile*, but it'll be a
	// function value, not a string, so it'll likely panic.
	Val js.Value
}

func (db *DirectiveBinding) Name() string        { return db.Val.Get("name").String() }
func (db *DirectiveBinding) Value() js.Value     { return db.Val.Get("value") }
func (db *DirectiveBinding) OldValue() js.Value  { return db.Val.Get("oldValue") }
func (db *DirectiveBinding) Expression() string  { return db.Val.Get("expression").String() }
func (db *DirectiveBinding) Arg() string         { return db.Val.Get("arg").String() }
func (db *DirectiveBinding) Modifiers() js.Value { return db.Val.Get("modifiers") }

// NewDirective creates a new directive.  It wraps js{Vue.directive}.
// https://vuejs.org/v2/api/#Vue-directive
func NewDirective(name string, opts ...DirectiveOption) *Directive {
	if len(opts) == 0 {
		// Retrieve the directive
		return &Directive{Value: js.Global().Get("Vue").Call("directive", name)}
	}
	c := &DirectiveConfig{Value: NewObject()}
	c.Option(opts...)
	if c.Short() != js.Undefined() {
		return &Directive{Value: js.Global().Get("Vue").Call("directive", name, c.Short())}
	}
	return &Directive{Value: js.Global().Get("Vue").Call("directive", name, c.Value)}
}

// Bind specifies the js{bind} directive hook function.  Called only once,
// when the directive is first bound to the element. This is where you can do
// one-time setup work.
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Bind(f func(el js.Value, binding *DirectiveBinding, vnode js.Value)) DirectiveOption {
	return makeDirectiveOption("bind", f)
}

// Inserted specifies the js{inserted} directive hook function.  Called when
// the bound element has been inserted into its parent node (this only
// guarantees parent node presence, not necessarily in-document).
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Inserted(f func(el js.Value, binding *DirectiveBinding, vnode js.Value)) DirectiveOption {
	return makeDirectiveOption("inserted", f)
}

// Update specifies the js{update} directive hook function.  Called after the
// containing component has updated, but possibly before its children have
// updated. The directive’s value may or may not have changed, but you can
// skip unnecessary updates by comparing the binding’s current and old values
// (see the Vue Guide on hook arguments).
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Update(f func(el js.Value, binding *DirectiveBinding, vnode, oldVnode js.Value)) DirectiveOption {
	return makeDirectiveUpdateOption("update", f)
}

// ComponentUpdated specifies the js{componentUpdated} directive hook
// function.  Called after the containing component and its children have
// updated.  https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func ComponentUpdated(f func(el js.Value, binding *DirectiveBinding, vnode, oldVode js.Value)) DirectiveOption {
	return makeDirectiveUpdateOption("componentUpdated", f)
}

// Unbind specifies the js{unbind} directive hook function.  Called only once,
// when the directive is unbound from the element.
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Unbind(f func(el js.Value, binding *DirectiveBinding, vnode js.Value)) DirectiveOption {
	return makeDirectiveOption("unbind", f)
}

func makeDirectiveOption(name string, f func(el js.Value, binding *DirectiveBinding, vnode js.Value)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.Set(name, js.NewCallback(
			func(thisNotSet js.Value, args []js.Value) interface{} {
				f(args[0],
					&DirectiveBinding{Val: args[1]},
					args[2])
				return nil
			}))
	}
}

func makeDirectiveUpdateOption(name string, f func(el js.Value, binding *DirectiveBinding, vnode, oldVnode js.Value)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.Set(name, js.NewCallback(
			func(thisNotSet js.Value, jsArgs []js.Value) interface{} {
				f(jsArgs[0],
					&DirectiveBinding{Val: jsArgs[1]},
					jsArgs[2],
					jsArgs[3])
				return nil
			}))
	}
}

// Short allows you to use the "function shorthand" style of directive
// definition, when you want the same behavior on bind and update, but don't
// care about the other hooks.  oldVnode is only used for the update hook; for
// the bind hook, it's nil.
// https://vuejs.org/v2/guide/custom-directive.html#Function-Shorthand
func Short(f func(el js.Value, binding *DirectiveBinding, vnode, oldVnode js.Value)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.SetShort(js.NewCallback(
			func(thisNotSet js.Value, jsArgs []js.Value) interface{} {
				var lastArg js.Value
				if len(jsArgs) == 4 {
					lastArg = jsArgs[3]
				}
				f(jsArgs[0],
					&DirectiveBinding{Val: jsArgs[1]},
					jsArgs[2],
					lastArg)
				return nil
			}))
	}
}

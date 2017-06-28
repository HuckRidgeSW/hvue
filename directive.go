package hvue

import (
	"github.com/gopherjs/gopherjs/js"
)

// Directive wraps a js{Vue.directive} object.
// https://vuejs.org/v2/api/#Vue-directive.
type Directive struct {
	*js.Object
}

// DirectiveBinding wraps the js{binding} slot of the directive hook argument.
// https://vuejs.org/v2/guide/custom-directive.html#Directive-Hook-Arguments
type DirectiveBinding struct {
	*js.Object
	Name       string      `js:"name"`
	Value      interface{} `js:"value"`
	OldValue   interface{} `js:"oldValue"`
	Expression string      `js:"expression"`
	Arg        string      `js:"arg"`
	Modifiers  *js.Object  `js:"modifiers"`
}

// NewDirective creates a new directive.  It wraps js{Vue.directive}.
// https://vuejs.org/v2/api/#Vue-directive
func NewDirective(name string, opts ...DirectiveOption) *Directive {
	if len(opts) == 0 {
		// Retrieve the directive
		return &Directive{Object: js.Global.Get("Vue").Call("directive", name)}
	}
	c := &DirectiveConfig{Object: NewObject()}
	c.Option(opts...)
	if c.Short != js.Undefined {
		return &Directive{Object: js.Global.Get("Vue").Call("directive", name, c.Short)}
	} else {
		return &Directive{Object: js.Global.Get("Vue").Call("directive", name, c.Object)}
	}
}

// Bind specifies the js{bind} directive hook function.  Called only once,
// when the directive is first bound to the element. This is where you can do
// one-time setup work.
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Bind(f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return makeDirectiveOption("bind", f)
}

// Inserted specifies the js{inserted} directive hook function.  Called when
// the bound element has been inserted into its parent node (this only
// guarantees parent node presence, not necessarily in-document).
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Inserted(f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return makeDirectiveOption("inserted", f)
}

// Update specifies the js{update} directive hook function.  Called after the
// containing component has updated, but possibly before its children have
// updated. The directive’s value may or may not have changed, but you can
// skip unnecessary updates by comparing the binding’s current and old values
// (see the Vue Guide on hook arguments).
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Update(f func(el *js.Object, binding *DirectiveBinding, vnode, oldVnode *js.Object)) DirectiveOption {
	return makeDirectiveUpdateOption("update", f)
}

// ComponentUpdated specifies the js{componentUpdated} directive hook
// function.  Called after the containing component and its children have
// updated.  https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func ComponentUpdated(f func(el *js.Object, binding *DirectiveBinding, vnode, oldVode *js.Object)) DirectiveOption {
	return makeDirectiveUpdateOption("componentUpdated", f)
}

// Unbind specifies the js{unbind} directive hook function.  Called only once,
// when the directive is unbound from the element.
// https://vuejs.org/v2/guide/custom-directive.html#Hook-Functions
func Unbind(f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return makeDirectiveOption("unbind", f)
}

func makeDirectiveOption(name string, f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.Object.Set(name, js.MakeFunc(
			func(thisNotSet *js.Object, jsArgs []*js.Object) interface{} {
				f(jsArgs[0],
					&DirectiveBinding{Object: jsArgs[1]},
					jsArgs[2])
				return nil
			}))
	}
}

func makeDirectiveUpdateOption(name string, f func(el *js.Object, binding *DirectiveBinding, vnode, oldVnode *js.Object)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.Object.Set(name, js.MakeFunc(
			func(thisNotSet *js.Object, jsArgs []*js.Object) interface{} {
				f(jsArgs[0],
					&DirectiveBinding{Object: jsArgs[1]},
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
func Short(f func(el *js.Object, binding *DirectiveBinding, vnode, oldVnode *js.Object)) DirectiveOption {
	return func(c *DirectiveConfig) {
		c.Short = js.MakeFunc(
			func(thisNotSet *js.Object, jsArgs []*js.Object) interface{} {
				var lastArg *js.Object = nil
				switch len(jsArgs) {
				case 3:
					// Do nothing
				case 4:
					lastArg = jsArgs[3]
				}

				f(jsArgs[0],
					&DirectiveBinding{Object: jsArgs[1]},
					jsArgs[2],
					lastArg)
				return nil
			})
	}
}

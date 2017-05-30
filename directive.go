package hvue

import (
	"github.com/gopherjs/gopherjs/js"
)

type Directive struct {
	*js.Object
}

type DirectiveBinding struct {
	*js.Object
	Name       string      `js:"name"`
	Value      interface{} `js:"value"`
	OldValue   interface{} `js:"oldValue"`
	Expression string      `js:"expression"`
	Arg        string      `js:"arg"`
	Modifiers  *js.Object  `js:"modifiers"`
}

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

func Bind(f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return makeDirectiveOption("bind", f)
}

func Inserted(f func(el *js.Object, binding *DirectiveBinding, vnode *js.Object)) DirectiveOption {
	return makeDirectiveOption("inserted", f)
}

func Update(f func(el *js.Object, binding *DirectiveBinding, vnode, oldVnode *js.Object)) DirectiveOption {
	return makeDirectiveUpdateOption("update", f)
}

func ComponentUpdated(f func(el *js.Object, binding *DirectiveBinding, vnode, oldVode *js.Object)) DirectiveOption {
	return makeDirectiveUpdateOption("componentUpdated", f)
}

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

package main

// From https://vuejs.org/v2/guide/custom-directive.html.  More direct links
// below.

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
)

func main() {
	go directiveHookArguments()
	go objectLiterals()
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/custom-directive.html#Directive-Hook-Arguments

func directiveHookArguments() {
	hvue.NewDirective("demo",
		hvue.Bind(func(el *js.Object, binding *hvue.DirectiveBinding, vnode *js.Object) {
			s := js.Global.Get("JSON").Get("stringify")
			el.Set("innerHTML",
				"name: "+s.Invoke(binding.Name).String()+"<br>"+
					"value: "+s.Invoke(binding.Value).String()+"<br>"+
					"expression: "+s.Invoke(binding.Expression).String()+"<br>"+
					"argument: "+s.Invoke(binding.Arg).String()+"<br>"+
					"modifiers: "+s.Invoke(binding.Modifiers).String()+"<br>"+
					"vnode keys: "+strings.Join(js.Keys(vnode), ", ")+"<br>"+
					"",
			)
		}))

	hvue.NewVM(
		hvue.El("#hook-arguments-example"),
		hvue.Data("message", "hello!"))
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/custom-directive.html#Object-Literals

func objectLiterals() {
	hvue.NewDirective("demo2",
		hvue.Short(func(el *js.Object, binding *hvue.DirectiveBinding, vnode, oldVnode *js.Object) {
			value := binding.Value.(map[string]interface{})
			fmt.Println(value["color"]) // => "white"
			fmt.Println(value["text"])  // => "hello!"
		}))
	hvue.NewVM(
		hvue.El("#object-literals"))
}

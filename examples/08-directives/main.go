package main

// From https://vuejs.org/v2/guide/custom-directive.html.  More direct links
// below.

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherwasm/js"
	"github.com/huckridgesw/hvue"
)

func main() {
	go directiveHookArguments()
	go objectLiterals()

	select {}
}

////////////////////////////////////////////////////////////////////////////////

// https://vuejs.org/v2/guide/custom-directive.html#Directive-Hook-Arguments

func directiveHookArguments() {
	hvue.NewDirective("demo",
		hvue.Bind(func(el js.Value, binding *hvue.DirectiveBinding, vnode js.Value) {
			stringify := func(i interface{}) string {
				return js.Global().Get("JSON").Call("stringify", i).String()
			}
			keysA := js.Global().Get("Object").Call("keys", vnode)
			var keys []string
			for i, l := 0, keysA.Length(); i < l; i++ {
				keys = append(keys, keysA.Index(i).String())
			}
			el.Set("innerHTML",
				"name: "+stringify(binding.Name())+"<br>"+
					"value: "+stringify(binding.Value())+"<br>"+
					"expression: "+stringify(binding.Expression())+"<br>"+
					"argument: "+stringify(binding.Arg())+"<br>"+
					"modifiers: "+stringify(binding.Modifiers())+"<br>"+
					"vnode keys: "+strings.Join(keys, ", ")+"<br>"+
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
		hvue.Short(func(el js.Value, binding *hvue.DirectiveBinding, vnode, oldVnode js.Value) {
			value := binding.Value()
			fmt.Println(value.Get("color").String()) // => "white"
			fmt.Println(value.Get("text").String())  // => "hello!"
		}))
	hvue.NewVM(
		hvue.El("#object-literals"))
}

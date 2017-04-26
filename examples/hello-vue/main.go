package main

import (
	"github.com/theclapp/hvue"
)

func main() {
	vm := hvue.NewVM(
		hvue.NewConfig(
			hvue.El("#app"),
			hvue.Data("message", "Hello, Vue!")))
	_ = vm
}

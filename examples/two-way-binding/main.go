package main

import "github.com/theclapp/hvue"

func main() {
	hvue.NewVM(
		hvue.El("#app-6"),
		hvue.Data("message", "Hello Vue!"))
}

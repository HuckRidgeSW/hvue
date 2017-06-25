package main

import "github.com/huckridge/hvue"

func main() {
	hvue.NewVM(
		hvue.El("#app-6"),
		hvue.Data("message", "Hello Vue!"))
}

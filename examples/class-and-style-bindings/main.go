package main

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

func main() {
	data := &struct {
		*js.Object
		IsActive bool `js:"isActive"`
	}{Object: hvue.NewObject()}
	data.IsActive = false

	hvue.NewVM(
		hvue.El("#object-syntax-1"),
		hvue.DataS(data))

	time.Sleep(2 * time.Second)
	data.IsActive = true
	println("isActive:", data.IsActive)
}

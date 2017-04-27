package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/theclapp/hvue"
)

type Data struct {
	*js.Object
	GroceryList []*ListItem `js:"groceryList"`
}

type ListItem struct {
	*js.Object
	Text string `js:"text"`
}

func main() {
	hvue.NewComponent("todo-item",
		hvue.Props("todo"),
		hvue.Template("<li>{{ todo.text }}</li>"))

	// This compiles and runs but wouldn't actually work in practice.  It'd be
	// nice to write a function that could take this and copy it into a new
	// structure with all the Object slots initialized correctly.
	testData := &Data{
		GroceryList: []*ListItem{
			&ListItem{Text: "stuff"},
		},
	}
	println("testData is", testData)

	hvue.NewVM(
		hvue.El("#app-7"),
		hvue.DataS(NewData(
			"Vegetables",
			"Cheese",
			"Whatever else humans are supposed to eat")))
}

func NewData(texts ...string) *Data {
	d := &Data{Object: hvue.NewObject()}
	d.GroceryList = NewGroceryList(texts)
	return d
}

func NewGroceryList(texts []string) []*ListItem {
	list := make([]*ListItem, len(texts))
	for i, text := range texts {
		list[i] = NewListItem(text)
	}
	return list
}

func NewListItem(Text string) *ListItem {
	item := &ListItem{Object: js.Global.Get("Object").New()}
	item.Text = Text
	return item
}

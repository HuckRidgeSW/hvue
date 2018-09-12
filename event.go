package hvue

import "github.com/gopherjs/gopherwasm/js"

// Event wraps the event object sent to v-on event handlers.
//
// I've only implemented enough of the Event type
// (https://developer.mozilla.org/en-US/docs/Web/API/Event) to implement
// example 07.
type Event struct {
	js.Value
}

type HTMLElement struct {
	js.Value
}

func (e *Event) Target() *HTMLElement {
	return &HTMLElement{Value: e.Get("target")}
}

func (et *HTMLElement) Select() {
	et.Call("select")
}

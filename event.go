package hvue

import "github.com/gopherjs/gopherjs/js"

// Event wraps the event object sent to v-on event handlers.
type Event struct {
	*js.Object
}

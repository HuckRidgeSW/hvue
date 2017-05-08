package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type Config struct {
	*js.Object
	El         string     `js:"el"`
	Data       *js.Object `js:"data"`
	Methods    *js.Object `js:"methods"`
	Props      *js.Object `js:"props"`
	Template   string     `js:"template"`
	Computed   *js.Object `js:"computed"`
	Components *js.Object `js:"components"`

	dataValue reflect.Value
}

type option func(*Config)

// Option sets the options specified.
func (c *Config) Option(opts ...option) {
	for _, opt := range opts {
		opt(c)
	}
}

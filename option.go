package hvue

import "github.com/gopherjs/gopherjs/js"

type Config struct {
	*js.Object
	El       string     `js:"el"`
	Data     *js.Object `js:"data"`
	Methods  *js.Object `js:"methods"`
	Props    *js.Object `js:"props"`
	Template string     `js:"template"`
	Computed *js.Object `js:"computed"`
}

type option func(*Config)

// Option sets the options specified.
func (c *Config) Option(opts ...option) {
	for _, opt := range opts {
		opt(c)
	}
}

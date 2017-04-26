package hvue

import "github.com/gopherjs/gopherjs/js"

type VM struct {
	*js.Object
}

type Config struct {
	*js.Object
	El   string     `js:"el"`
	Data *js.Object `js:"data"`
	// Template string     `js:"template"`
}

type option func(*Config)

func NewVM(opts ...option) *VM {
	c := &Config{Object: js.Global.Get("Object").New()}
	c.Data = js.Global.Get("Object").New()
	c.Option(opts...)
	return &VM{Object: js.Global.Get("Vue").New(c)}
}

func El(selector string) option {
	return func(c *Config) {
		c.El = selector
	}
}

func Data(name, value string) option {
	return func(c *Config) {
		c.Data.Set(name, value)
	}
}

// func Template(t string) option {
// 	return func(c *Config) {
// 		c.Template = t
// 	}
// }

// Option sets the options specified.
func (f *Config) Option(opts ...option) {
	for _, opt := range opts {
		opt(f)
	}
}
